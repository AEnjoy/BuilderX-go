package v1

type ApiGroup struct {
	TaskApi
	FileApi
}

var ApiGroupApp = new(ApiGroup)
