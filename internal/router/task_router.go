package router

import (
	"github.com/gin-gonic/gin"
	"github.com/issueye/log-agent/internal/controller"
)

type JobRouter struct {
	Name    string
	control *controller.TaskController
}

func NewJobRouter() *JobRouter {
	return &JobRouter{
		Name:    "cron",
		control: controller.NewTaskController(),
	}
}

func (job JobRouter) Register(group *gin.RouterGroup) {
	f := group.Group(job.Name)
	f.GET("list", job.control.List)
	f.POST("create", job.control.Create)
	f.PUT("modify/:id", job.control.Modify)
	f.PUT("modifyStatus/:id", job.control.ModifyStatus)
	f.DELETE("batchDelete", job.control.BatchDelete)
	f.DELETE("delete/:id", job.control.Delete)
}
