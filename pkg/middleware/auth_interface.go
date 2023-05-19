package middleware

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AuthInterface interface {
	// 获取jwt标识
	GetJwtRealm() string
	// jwt 秘钥
	GetJwtKey() string
	// 超时
	GetJwtTimeOut() int64
	// 刷新时间
	GetJwtMaxRefresh() int64
	// 有效载荷处理
	PayloadFunc(data interface{}) jwt.MapClaims
	// 解析Claims
	IdentityHandler(c *gin.Context) interface{}
	// 用户登录
	Login(c *gin.Context) (interface{}, error)
	// 用户登录校验成功处理
	Authorizator(data interface{}, c *gin.Context) bool
	// 用户登录校验失败处理
	Unauthorized(ctx *gin.Context, code int, message string)
	// 登录成功后的响应
	LoginResponse(ctx *gin.Context, _ int, token string, expires time.Time)
	// 用户登出
	LogoutResponse(ctx *gin.Context, _ int)
	// 刷新token
	RefreshResponse(ctx *gin.Context, _ int, token string, expires time.Time)
}
