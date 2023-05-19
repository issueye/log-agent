package initialize

import (
	"fmt"
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/issueye/log-agent/internal/config"
	"github.com/issueye/log-agent/internal/global"
	"github.com/issueye/log-agent/internal/router"
	"github.com/issueye/log-agent/pkg/middleware"
	orange_validator "github.com/issueye/log-agent/pkg/validator"
)

func InitServer() {
	mode := config.GetParam("SERVER-MODE", "debug").String()
	gin.SetMode(mode)
	// gin引擎对象
	global.Router = gin.New()
	// 注册一个form表单验证器
	orange_validator.RegisterValidator()

	// 加载中间件
	global.Router.Use(middleware.CORSMiddleware([]string{}))       // 处理前端跨域
	global.Router.Use(middleware.GinLogger(global.Logger))         // 日志记录
	global.Router.Use(middleware.GinRecovery(global.Logger, true)) // 服务恐慌处理

	// 设置一个静态文件服务器
	global.Router.Static("/www", "./runtime/static")

	// 告诉服务文件的MIME类型
	_ = mime.AddExtensionType(".js", "application/javascript")
	_ = mime.AddExtensionType(".css", "text/css")
	_ = mime.AddExtensionType(".woff", "application/font-woff")
	_ = mime.AddExtensionType(".woff2", "application/font-woff2")
	_ = mime.AddExtensionType(".ttf", "application/font-ttf")
	_ = mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")

	// 注册路由
	router.InitRouter(global.Router)
	// 端口号为命令行提供
	port := config.GetParam("SERVER-PORT", "10061").Int()
	global.HttpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: global.Router,
	}
}
