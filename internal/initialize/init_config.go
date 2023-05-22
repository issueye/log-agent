package initialize

import "github.com/issueye/log-agent/internal/config"

func InitConfig() {
	config.SetParamExist("SERVER-PORT", "10061", "端口号")
	config.SetParamExist("SERVER-MODE", "debug", `服务运行模式， debug\release`)

	config.SetParamExist("LOG-MAX-SIZE", "10", "日志大小")
	config.SetParamExist("LOG-MAX-BACKUPS", "10", "最大备份数")
	config.SetParamExist("LOG-MAX-AGE", "10", "保存天数")
	config.SetParamExist("LOG-COMPRESS", "true", "是否压缩")
	config.SetParamExist("LOG-LEVEL", "-1", "日志输出等级")

	config.SetParamExist("LOG-AGENT-ASYNC", "false", "日志同步模式 true 异步 false 同步")
}
