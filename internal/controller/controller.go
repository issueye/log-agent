package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/issueye/log-agent/pkg/res"
)

type Controller struct {
	*gin.Context
}

func New(ctx *gin.Context) *Controller {
	return &Controller{
		Context: ctx,
	}
}

func (c *Controller) Success() {
	res.Success(c.Context)
}

func (c *Controller) SuccessByMsg(msg string) {
	res.SuccessByMsg(c.Context, msg)
}

func (c *Controller) SuccessByMsgf(fmtStr string, args ...any) {
	res.SuccessByMsgf(c.Context, fmtStr, args...)
}

func (c *Controller) SuccessData(data interface{}) {
	res.SuccessData(c.Context, data)
}

func (c *Controller) SuccessAutoData(req interface{}, data interface{}) {
	res.SuccessAutoData(c.Context, req, data)
}

func (c *Controller) SuccessPage(data interface{}) {
	res.SuccessPage(c.Context, data)
}

func (c *Controller) Fail() {
	res.Fail(c.Context, res.BAD_REQUEST)
}

func (c *Controller) FailByMsg(msg string) {
	res.FailByMsg(c.Context, msg)
}

func (c *Controller) FailByMsgf(fmtStr string, args ...any) {
	res.FailByMsgf(c.Context, fmtStr, args...)
}

func (c *Controller) FailBind(err error) {
	res.FailBind(c.Context, err)
}
