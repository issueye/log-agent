package model

type Page struct {
	PageNum  int64 `json:"pageNum" form:"pageNum"`   // 页数
	PageSize int64 `json:"pageSize" form:"pageSize"` // 页码
	Total    int64 `json:"total"`                    // 总数  由服务器返回回去
}
