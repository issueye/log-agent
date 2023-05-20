package model

type Monitor struct {
	ID         string `gorm:"column:id;primaryKey;type:nvarchar(100);" json:"id"`       // 编码
	Name       string `gorm:"column:name;type:nvarchar(100);" json:"name"`              // 名称
	LogPath    string `gorm:"column:log_path;type:nvarchar(400);" json:"logPath"`       // 日志路径
	Level      int    `gorm:"column:level;type:int;" json:"level"`                      // 日志等级 如果没有设置则默认全部类型  -1 debug 0 info 1 warn 2 error
	ScriptPath string `gorm:"column:script_path;type:nvarchar(400);" json:"scriptPath"` // 脚本路径 可为空
	State      bool   `gorm:"column:state;type:bit;" json:"state"`                      // 是否启用
	CreateTime string `gorm:"column:create_time;type:nvarchar(100);" json:"createTime"` // 创建时间
}

func (Monitor) TableName() string {
	return "monitor"
}

type CreateMonitor struct {
	Name       string `json:"name" binding:"required" label:"名称"`         // 名称
	LogPath    string `json:"logPath" binding:"required" label:"日志路径"`    // 日志路径
	Level      int    `json:"level" binding:"required" label:"日志等级"`      // 日志等级 如果没有设置则默认全部类型
	ScriptPath string `json:"scriptPath" binding:"required" label:"脚本路径"` // 脚本路径 可为空
}

type ModifyMonitor struct {
	ID         string `json:"id" binding:"required" label:"编码"` // 编码
	Name       string `json:"name"`                             // 名称
	LogPath    string `json:"logPath"`                          // 日志路径
	Level      int    `json:"level"`                            // 日志等级 如果没有设置则默认全部类型
	ScriptPath string `json:"scriptPath"`                       // 脚本路径 可为空
}

// 查询条件
type QueryMonitor struct {
	Condition string `json:"condition" form:"condition"` // 条件
	Page
}
