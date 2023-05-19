package initialize

import (
	"github.com/issueye/log-agent/internal/global"
	"github.com/issueye/log-agent/internal/service"
)

// 初始化其他数据
func InitData() {
	err := service.NewUser(global.DB).CreateAdminNonExistent()
	if err != nil {
		panic("初始化数据失败，失败原因：" + err.Error())
	}

	// 初始化定时任务
	service.NewTask(global.DB).CreateJobNoExistent()
}
