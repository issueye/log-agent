package model

// Task
// 定时任务
type Task struct {
	ID         string `gorm:"column:id;primaryKey;type:nvarchar(100);" json:"id"`     // 任务编码
	Name       string `gorm:"column:name;type:nvarchar(100);" json:"name"`            // 任务名称
	Expression string `gorm:"column:expression;type:nvarchar(50);" json:"expression"` // 时间表达式
	State      bool   `gorm:"column:state;type:bit;" json:"state"`                    // 是否启用
	Path       string `gorm:"column:path;nvarchar(200);" json:"path"`                 // 脚本路径
	Mark       string `gorm:"column:mark;type:nvarchar(500);" json:"mark"`            // 备注
}

func (Task) TableName() string {
	return "task_info"
}

type TaskNoticeType int64

const (
	DEL_JOB TaskNoticeType = iota
	ADD_JOB
	MODIFY_JOB
)

// NoticeJob
// 通知
type NoticeJob struct {
	Task
	Type TaskNoticeType
}

type CreateTask struct {
	Name       string `json:"name" binding:"required" label:"任务名称"`        // 任务名称
	Expression string `json:"expression" binding:"required" label:"时间表达式"` // 时间表达式
	Path       string `json:"path" binding:"required" label:"脚本路径"`        // 脚本路径
	Mark       string `json:"mark"`                                        // 备注
}

type DelTask struct {
	Ids []string `json:"ids"` // 定时任务编码
}

type ModifyTask struct {
	ID         string `json:"id" binding:"required" label:"任务编码"` // 编码
	Name       string `json:"name"`                               // 名称
	Expression string `json:"expression"`                         // 时间表达式
	State      bool   `json:"state"`                              // 是否启用
	Path       string `json:"path"`                               // 脚本路径
	Mark       string `json:"mark"`                               // 备注
}

type QueryTask struct {
	Name string `json:"name" form:"name"` // 任务名称
	Mark string `json:"mark" form:"mark"` // 任务描述
	Page
}
