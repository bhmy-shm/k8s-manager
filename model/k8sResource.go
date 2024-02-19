package model

import (
	"context"
	"fmt"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/gormx"
	gofkHash "github.com/bhmy-shm/gofks/core/utils/hash"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
	"time"
)

var _ K8sResourceModel = (*defaultResource)(nil)

type (
	K8sResourceModel interface {
		resQueryModel
		resEditModel
	}

	resQueryModel interface {
		Trans(context.Context, gormx.TransFunc) error
		Query(context.Context, interface{}, ...gormx.SqlOptions) error
		QueryCount(context.Context, *int64, ...gormx.SqlOptions) error
	}
	resEditModel interface {
		AutoMigrates()
		Insert(context.Context, *K8sResource) error
		Update(context.Context, *K8sResource) error
		Delete(context.Context, string) error
	}

	K8sResource struct {
		gorm.Model
		obj             *unstructured.Unstructured `gorm:"-"`
		objBytes        []byte                     `gorm:"-"`
		NameSpace       string                     `gorm:"column:namespace"`
		Name            string                     `gorm:"column:name"`
		ResourceVersion string                     `gorm:"column:resource_version"`
		Hash            string                     `gorm:"column:hash;unique;"`

		//gvr 和kind相关
		Group    string `gorm:"column:group"`
		Version  string `gorm:"column:version"`
		Resource string `gorm:"column:resource"`
		Kind     string `gorm:"column:kind"`

		// uid是唯一的 。owner 不一定有
		Owner  string `gorm:"column:owner"`
		Uid    string `gorm:"column:uid;unique;index"`
		Object string `gorm:"-"` //yaml对象 `gorm:"column:object"`
	}

	defaultResource struct {
		session   gormx.SqlSession
		tableName string
	}
)

// ----------- 初始化

func NewK8sResourceModel(c *gofkConf.Config) K8sResourceModel {
	if c.GetDB().IsLoad() {
		model := newK8sModel(c)
		model.AutoMigrates()
		return model
	}
	return nil
}

func newK8sModel(c *gofkConf.Config) *defaultResource {
	return &defaultResource{
		session:   gormx.NewSql(c.GetDB()),
		tableName: "k8sResource",
	}
}

func (m *defaultResource) AutoMigrates() {
	err := m.session.RawDB().AutoMigrate(&K8sResource{})
	if err != nil {
		panic(err)
	}
}

func (m *defaultResource) TableName() string {
	return m.tableName
}

func (m *defaultResource) Trans(ctx context.Context, transFunc gormx.TransFunc) error {
	return m.session.TransactCtx(ctx, transFunc)
}

func (m *defaultResource) Query(ctx context.Context, result interface{}, opts ...gormx.SqlOptions) error {
	db := m.session.RawDB().Model(&K8sResource{})
	for _, opt := range opts {
		opt(db)
	}
	return m.session.QueryFromDB(ctx, db, result)
}

func (m *defaultResource) QueryCount(ctx context.Context, total *int64, opts ...gormx.SqlOptions) error {
	db := m.session.RawDB().Model(K8sResource{})
	for _, opt := range opts {
		opt(db)
	}
	return m.session.CountFromDB(ctx, db, total)
}

func (m *defaultResource) Insert(ctx context.Context, data *K8sResource) error {
	return m.session.TransactCtx(ctx, func(context context.Context, tx *gorm.DB) error {

		resDB := tx.Model(&K8sResource{})

		resDB = data.ClausesOnConflict(resDB) //解决资源反复插入的冲突问题

		return resDB.Create(data).Error
	})
}

func (m *defaultResource) Update(ctx context.Context, data *K8sResource) error {
	return m.session.TransactCtx(ctx, func(context context.Context, tx *gorm.DB) error {

		data.prepare() //资源准备工作

		resDB := tx.Model(&K8sResource{})

		return resDB.
			Where("uid = ?", data.Uid).
			Where("hash != ?", data.Hash).
			Updates(data).Error
	})
}

func (m *defaultResource) Delete(ctx context.Context, uid string) error {
	return m.session.TransactCtx(ctx, func(context context.Context, tx *gorm.DB) error {
		return tx.Model(&K8sResource{}).Delete(&K8sResource{}).Where("uid=", uid).Error
	})
}
func (m *K8sResource) TableName() string {
	return "k8sResource"
}

// ------- 数据处理方法 --------

func NewPureResource() *K8sResource {
	return &K8sResource{}
}

func NewResource(obj runtime.Object, restmapper meta.RESTMapper) (*K8sResource, error) {
	o := &unstructured.Unstructured{}

	b, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, o)
	if err != nil {
		return nil, err
	}

	//获取gvk
	gvk := o.GroupVersionKind()
	mapping, err := restmapper.RESTMapping(gvk.GroupKind())
	if err != nil {
		return nil, err
	}

	// 赋值
	retObj := &K8sResource{obj: o, objBytes: b}
	retObj.Model.CreatedAt = o.GetCreationTimestamp().Time
	retObj.Uid = string(o.GetUID())
	retObj.Name = o.GetName()
	retObj.NameSpace = o.GetNamespace()
	retObj.Group = gvk.Group
	retObj.Version = gvk.Version
	retObj.Kind = gvk.Kind
	retObj.Resource = mapping.Resource.Resource
	retObj.ResourceVersion = o.GetResourceVersion()

	return retObj, nil
}

func (r *K8sResource) prepare() {
	if len(r.obj.GetOwnerReferences()) > 0 {
		r.Owner = string(r.obj.GetOwnerReferences()[0].UID)
	}
	r.Hash = gofkHash.Md5Hex(r.objBytes) //保存 资源的md5值
	objectJson, err := yaml.YAMLToJSON(r.objBytes)
	if err != nil {
		objectJson = []byte(fmt.Errorf("obj Yaml To json failed: %v", err).Error())
	}
	r.Object = string(objectJson)
}

func (r *K8sResource) ClausesOnConflict(db *gorm.DB) *gorm.DB {

	//prepare 主要构建 object（pod-yaml）对象
	r.prepare()

	//这里的 clause.Column 必须是唯一条件约束字段
	//ALTER TABLE "k8sResource" ADD CONSTRAINT unique_uid UNIQUE (uid);
	return db.Clauses(clause.OnConflict{

		Columns: []clause.Column{{Name: "uid"}},

		DoUpdates: clause.Assignments(
			map[string]interface{}{
				"resource_version": r.ResourceVersion,
				"updated_at":       time.Now(),
			},
		),
	})
}

func WithNamespace(ns string) gormx.SqlOptions {
	return func(newTx *gorm.DB) *gorm.DB {
		if len(ns) > 0 {
			return newTx.Where("namespace = ?", ns)
		}
		return newTx
	}
}

func WithGroup(gp string) gormx.SqlOptions {
	return func(newTx *gorm.DB) *gorm.DB {
		if len(gp) > 0 {
			return newTx.Where("group = ?", gp)
		}
		return newTx
	}
}

func WithVersion(version string) gormx.SqlOptions {
	return func(newTx *gorm.DB) *gorm.DB {
		if len(version) > 0 {
			return newTx.Where("version = ?", version)
		}
		return newTx
	}
}

func WithResource(rs string) gormx.SqlOptions {
	return func(newTx *gorm.DB) *gorm.DB {
		if len(rs) > 0 {
			return newTx.Where("resource = ?", rs)
		}
		return newTx
	}
}
