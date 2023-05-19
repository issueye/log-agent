package config

// Param
// 参数
type Param struct {
	ID    int64  `gorm:"column:id;primaryKey;autoIncrement:false;type:int" json:"id"` // 编码
	Name  string `gorm:"column:name;type:nvarchar(50)" json:"name"`                   // 参数名称
	Value string `gorm:"column:value;type:nvarchar(50)" json:"value"`                 // 参数值
	Mark  string `gorm:"column:mark;type:nvarchar(500)" json:"mark"`                  // 备注
}

func (Param) TableName() string {
	return "param_info"
}
