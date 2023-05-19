package model

// User
// 用户信息
type User struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement:false;type:int" json:"id"` // 编码
	Account  string `gorm:"column:account;type:nvarchar(50)" json:"account"`             // uid 登录名
	Name     string `gorm:"column:name;type:nvarchar(50)" json:"name"`                   // 用户姓名
	Password string `gorm:"column:password;type:nvarchar(50)" json:"password"`           // 密码
	Mark     string `gorm:"column:mark;type:nvarchar(500)" json:"mark"`                  // 备注
	State    int    `gorm:"column:state;type:int" json:"state"`                          // 状态 0 停用 1 启用
}

type CreateUser struct {
	Account  string `json:"account"`  // uid 登录名
	Name     string `json:"name"`     // 用户姓名
	Password string `json:"password"` // 密码
	Mark     string `json:"mark"`     // 备注
}

type ModifyUser struct {
	ID       int64  `json:"id"`       // 编码
	Account  string `json:"account"`  // uid 登录名
	Name     string `json:"name"`     // 用户姓名
	Password string `json:"password"` // 密码
	Mark     string `json:"mark"`     // 备注
}

type StatusUser struct {
	ID    int64 `json:"id"`    // 编码
	State bool  `json:"state"` // 备注
}

type QueryUser struct {
	Account string `json:"account"` // uid 登录名
	Name    string `json:"name"`    // 用户姓名
	Mark    string `json:"mark"`    // 备注
}

// Login
// 用户登录
type Login struct {
	Account  string `json:"account"`  // 登录名
	Password string `json:"password"` // 密码
}

func (User) TableName() string {
	return "user_info"
}
