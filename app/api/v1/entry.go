package v1

import "github.com/gin-gonic/gin"

type TaskApi struct{}
type ApiGroup struct {
	TaskApi
}

var ApiGroupApp = new(ApiGroup)

func (b *TaskApi) CreateTask(c *gin.Context) {

}
