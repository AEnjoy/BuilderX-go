package v1

type ApiGroup struct {
	TaskApi
	FileApi
	BaseApi
}

var ApiGroupApp = new(ApiGroup)
