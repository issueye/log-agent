package router

import (
	"github.com/gin-gonic/gin"
	"github.com/issueye/log-agent/internal/config"
	"github.com/issueye/log-agent/internal/controller"
	"github.com/issueye/log-agent/internal/global"
	"github.com/issueye/log-agent/internal/middleware"
)

type IRouters interface {
	Register(group *gin.RouterGroup)
}

func InitRouter(r *gin.Engine) {
	name := config.GetParam("api-name", "api").String()
	v := config.GetParam("api-version", "v1").String()

	apiName := r.Group(name)
	version := apiName.Group(v)
	global.Auth = middleware.NewAuth()

	// 用户鉴权
	version.POST("login", global.Auth.LoginHandler)
	version.GET("logout", global.Auth.LogoutHandler)
	version.GET("refreshToken", global.Auth.RefreshHandler)

	// 鉴权
	version.Use(global.Auth.MiddlewareFunc())
	registerVersionRouter(version,
		NewParamsRouter(), // 参数配置
		NewJobRouter(),    // 定时任务
	)

	// monitor
	monitor := version.Group("monitor", global.Auth.MiddlewareFunc())
	{
		monitor.GET("", controller.ListMonitor)
		monitor.POST("", controller.AddMonitor)
		monitor.PUT("", controller.ModifyMonitor)
		monitor.PUT("state/:id", controller.ModifyStateMonitor)
		monitor.DELETE("/:id", controller.DelMonitor)

		// websocket
		r.GET("ws/monitor/:id", controller.WsLogView)
	}
}

// registerRouter 注册路由
func registerVersionRouter(r *gin.RouterGroup, iRouters ...IRouters) {
	for _, iRouter := range iRouters {
		iRouter.Register(r)
	}
}
