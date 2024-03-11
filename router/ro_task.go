package router

import (
	v1 "github.com/aenjoy/BuilderX-go/app/api/v1"
	"github.com/gin-gonic/gin"
)

type TaskRouter struct{}

func (s *TaskRouter) InitRouter(Router *gin.RouterGroup) {
	taskRouter := Router.Group("task")
	taskApi := v1.ApiGroupApp.TaskApi
	taskRouter.POST("/create", taskApi.CreateTask)
}
