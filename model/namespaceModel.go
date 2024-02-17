package model

type (
	NamespaceModel struct {
		Total int64      `json:"total"`
		List  []*NsModel `json:"list"`
	}
	NsModel struct {
		Name string `json:"name"`
	}
)
