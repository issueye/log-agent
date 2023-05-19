package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/issueye/log-agent/internal/global"
	"github.com/issueye/log-agent/internal/model"
	"github.com/issueye/log-agent/internal/service"
)

// 添加一条监听
func AddMonitor(ctx *gin.Context) {
	control := New(ctx)
	req := new(model.CreateMonitor)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	err = service.NewMonitor(global.DB).Create(req)
	if err != nil {
		control.FailByMsgf("添加监听信息失败，失败原因：%s", err.Error())
		return
	}

	control.Success()
}

// 修改监听信息
func ModifyMonitor(ctx *gin.Context) {
	control := New(ctx)
	req := new(model.ModifyMonitor)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	err = service.NewMonitor(global.DB).Modify(req)
	if err != nil {
		control.FailByMsgf("修改监听信息失败，失败原因：%s", err.Error())
		return
	}

	control.Success()
}

// 停用启用监听
func ModifyStateMonitor(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		control.FailByMsg("传入参数[id]不能为空")
		return
	}

	state, err := service.NewMonitor(global.DB).ModifyState(id)
	if err != nil {
		control.FailByMsgf("修改监听信息失败，失败原因：%s", err.Error())
		return
	}

	if state {
		// TODO 开启监听
	} else {
		// TODO 关闭监听
	}

	control.Success()
}
