package handler

import (
	"context"
	"fmt"
	"github.com/bhmy-shm/gofks/core/logx"
	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"log"
	"manager/internal/maps"
	"manager/model"
	"manager/svc"
)

var RateLimitQue workqueue.RateLimitingInterface

func init() {
	RateLimitQue = workqueue.NewRateLimitingQueue(&workqueue.BucketRateLimiter{
		Limiter: rate.NewLimiter(2, 1),
	})
}

type RateLimitResource struct {
	Type     string
	Resource interface{}
}

func NewRateLimitResource(t string, obj interface{}) *RateLimitResource {
	return &RateLimitResource{Type: t, Resource: obj}
}

type ResourceHandler struct {
	MetaResMapper  *maps.MetaResMapper `inject:"-"`
	ServiceContext *svc.ServiceContext `inject:"-"`
}

func (r *ResourceHandler) OnAdd(obj interface{}) {

	if o, ok := obj.(runtime.Object); ok {
		ret, err := model.NewResource(o, r.MetaResMapper.Mapper)
		if err != nil {
			logx.Errorf("OnAdd Resource解析失败：%v", err)
			return
		}

		//执行入库
		err = r.ServiceContext.ResourceModel.Insert(context.Background(), ret)
		if err != nil {
			logx.Errorf("OnAdd Resource入库失败：%v", err)
		}
	}
}

func (r *ResourceHandler) OnUpdate(oldObj, newObj interface{}) {
	if o, ok := newObj.(runtime.Object); ok {

		ret, err := model.NewResource(o, r.MetaResMapper.Mapper)
		if err != nil {
			logx.Errorf("OnAdd Resource解析失败：%v", err)
			return
		}

		err = r.ServiceContext.ResourceModel.Update(context.Background(), ret)
		if err != nil {
			//todo 插入出错。  出错后怎么办: 1、打印 2、记录日志 3、放到队列里，再处理
			logx.Errorf("OnAdd Resource 更新失败：%v", err)
		}
	}
}

func (r *ResourceHandler) OnDelete(obj interface{}) {

	getObj, err := meta.Accessor(obj)
	if err != nil {
		logx.Errorf("OnDelete Resource 解析失败：%v", err)
		return
	}

	err = r.ServiceContext.ResourceModel.Delete(context.Background(), string(getObj.GetUID()))
	if err != nil {
		logx.Errorf("OnDelete Resource 删除失败：%v", err)
	}
}

func (r *ResourceHandler) RateLimitConsumer() {
	log.Println("依赖注入运行消费者")
	go func() {
		for {
			get_obj, _ := RateLimitQue.Get()
			if rm_obj, ok := get_obj.(*RateLimitResource); ok {
				fmt.Print("类型是:", rm_obj.Type)
				obj, err := meta.Accessor(rm_obj.Resource)
				if err == nil {
					fmt.Printf(" 资源是:%s/%s \n", obj.GetNamespace(),
						obj.GetName())
				}
			}
		}
	}()
}
