package model

type ReqParamList struct {
	Condition string `json:"condition" form:"condition"` // 条件
	Page
}
