package router

import (
	"github.com/gin-gonic/gin"
	"github.com/issueye/log-agent/internal/controller"
)

var params *ParamsRouter

type ParamsRouter struct {
	Name    string
	control *controller.ParamsController
}

func NewParamsRouter() *ParamsRouter {
	return &ParamsRouter{
		Name:    "param",
		control: controller.NewParamsController(),
	}
}

func (c ParamsRouter) Register(group *gin.RouterGroup) {
	r := group.Group(c.Name)
	r.GET("", c.control.GetParamList)
	r.GET("/:name", c.control.GetParamById)
	r.PUT("", c.control.ModifyParam)
}
