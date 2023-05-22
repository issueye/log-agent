package initialize

import (
	"fmt"

	"github.com/issueye/log-agent/internal/agent"
	"github.com/issueye/log-agent/internal/global"
	"github.com/issueye/log-agent/internal/model"
	"github.com/issueye/log-agent/internal/service"
)

func InitMonitor() {
	req := new(model.QueryMonitor)
	list, err := service.NewMonitor(global.DB).Query(req)
	if err != nil {
		global.Log.Errorf("查询监听文件信息失败，失败原因：%s", err.Error())
		return
	}

	// 开启启用的任务
	for _, m := range list {
		if m.State {
			a, err := agent.New(m.ID, m.LogPath, m.ScriptPath, m.Level)
			if err != nil {
				fmt.Printf("【%s】监听失败，失败原因：%s", a.Path, err.Error())
				continue
			}
			a.Listen()
		}
	}
}
