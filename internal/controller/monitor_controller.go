package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/issueye/log-agent/internal/agent"
	"github.com/issueye/log-agent/internal/global"
	"github.com/issueye/log-agent/internal/model"
	"github.com/issueye/log-agent/internal/service"
	"github.com/issueye/log-agent/pkg/ws"
)

// websocket 升级并跨域
var (
	upgrade = &websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 日志查看 ws
func WsLogView(ctx *gin.Context) {
	control := New(ctx)
	id := control.Param("id")
	if id == "" {
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	// 升级为 websocket
	conn, err := upgrade.Upgrade(control.Writer, control.Request, nil)
	if err != nil {
		control.FailByMsgf("升级协议失败，失败原因：%s", err.Error())
	}

	ws.NewConn(id, conn)
	control.Success()
}

// 获取监听列表
func ListMonitor(ctx *gin.Context) {
	control := New(ctx)
	req := new(model.QueryMonitor)
	err := control.Bind(req)
	if err != nil {
		control.FailBind(err)
		return
	}

	list, err := service.NewMonitor(global.DB).Query(req)
	if err != nil {
		control.FailByMsgf("查询失败，失败原因：%s", err.Error())
		return
	}

	control.SuccessAutoData(req, list)
}

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

	// 如果启用的状态不允许修改
	data, err := service.NewMonitor(global.DB).GetById(req.ID)
	if err != nil {
		control.FailByMsgf("查询信息失败，失败原因：%s", err.Error())
		return
	}

	if data.State {
		control.FailByMsg("已启用，不允许修改")
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

	// 先查询是否存在
	data, err := service.NewMonitor(global.DB).GetById(id)
	if err != nil {
		control.FailByMsgf("查询信息失败，失败原因：%s", err.Error())
		return
	}

	state, err := service.NewMonitor(global.DB).ModifyState(id)
	if err != nil {
		control.FailByMsgf("修改监听信息失败，失败原因：%s", err.Error())
		return
	}

	if state {
		a, err := agent.New(data.ID, data.LogPath, data.ScriptPath, data.Level)
		if err != nil {
			control.FailByMsgf("创建监听失败，失败原因：%s", err.Error())
			return
		}
		a.Listen()
	} else {
		agent.Del(data.ID)
	}

	control.Success()
}

// 删除日志监听
func DelMonitor(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		control.FailByMsg("传入参数[id]不能为空")
		return
	}

	// 判断监听是否还在进行中，如果还在进行中，则不允许删除
	m, err := service.NewMonitor(global.DB).GetById(id)
	if err != nil {
		control.FailByMsgf("查询信息失败，失败原因：%s", err.Error())
		return
	}

	if m.State {
		control.FailByMsg("监听还在进行中，请先关闭再删除")
		return
	}

	err = service.NewMonitor(global.DB).DelMonitor(id)
	if err != nil {
		control.FailByMsgf("删除日志监听失败，失败原因：%s", err.Error())
		return
	}

	control.Success()
}
