package global

import (
	"embed"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/issueye/log-agent/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	Log         *zap.SugaredLogger
	Logger      *zap.Logger
	Router      *gin.Engine
	HttpServer  *http.Server
	Stop        = make(chan struct{})
	SwaggerJson []byte
	Auth        *jwt.GinJWTMiddleware
	Asset       embed.FS
)

var JobChan = make(chan *model.NoticeJob, 100)
