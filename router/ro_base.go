package router

import (
	v1 "github.com/aenjoy/BuilderX-go/app/api/v1"
	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	baseApi := v1.ApiGroupApp.BaseApi
	baseRouter.GET("/get-token", baseApi.GetToken)
	baseRouter.POST("/reset-token", baseApi.ResetCookie)
	baseRouter.POST("/upload", baseApi.UploadFiles) //body file:file path:path task:task
}
