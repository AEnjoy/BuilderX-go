package router

import (
	v1 "github.com/aenjoy/BuilderX-go/app/api/v1"
	"github.com/gin-gonic/gin"
)

type FileRouter struct{}

func (s *FileRouter) InitRouter(Router *gin.RouterGroup) {
	taskRouter := Router.Group("file")
	taskApi := v1.ApiGroupApp.FileApi
	taskRouter.POST("/list", taskApi.DirList)
	taskRouter.POST("/select", taskApi.FileSelect)
}
