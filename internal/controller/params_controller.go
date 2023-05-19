package controller

import (
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/issueye/log-agent/internal/config"
	"github.com/issueye/log-agent/internal/model"
)

type ParamsController struct {
	Controller
}

func NewParamsController() *ParamsController {
	c := new(ParamsController)
	return c
}

// GetParamList doc
//
// @tags        服务参数管理
// @Summary     获取参数列表
// @Description 获取参数列表
// @Produce     json
// @Param       isNotPaging query    string                            false "是否需要分页， 默认需要， 如果不分页 传 true"
// @Param       pageNum     query    string                            true  "页码， 如果不分页 传 0"
// @Param       pageSize    query    string                            true  "一页大小， 如果不分页 传 0"
// @Param       name        query    string                            false "参数名称"
// @Param       status      query    int                               false "参数状态"
// @Param       secens      query    string                            false "参数域"
// @Success     200         {object} res.Full{data=[]config.ParamInfo} true  "code: 200 成功"
// @Failure     500         {object} res.Base                          true  "错误返回内容"
// @Router      /granada/api/v1/param [get]
// @Security    ApiKeyAuth
func (ParamsController) GetParamList(ctx *gin.Context) {
	control := New(ctx)

	info := new(model.ReqParamList)
	err := control.BindQuery(info)
	if err != nil {
		control.FailByMsgf("绑定参数失败，失败原因：%s", err.Error())
		return
	}

	list, err := config.GetParamList(info.Condition)
	if err != nil {
		control.FailByMsgf("获取参数列表失败，失败原因：%s", err.Error())
		return
	}

	info.Total = int64(len(list))
	control.SuccessAutoData(info, list)
}

// GetParamById doc
//
// @tags        服务参数管理
// @Summary     通过参数名称获取参数
// @Description 通过参数名称获取参数
// @Produce     json
// @Param       name path     string                        false "参数名称"
// @Success     200  {object} res.Full{data=model.ResParam} true  "code: 200 成功"
// @Failure     500  {object} res.Base                      true  "错误返回内容"
// @Router      /granada/api/v1/param/{name} [get]
// @Security    ApiKeyAuth
func (ParamsController) GetParamById(ctx *gin.Context) {
	control := New(ctx)

	name := control.Param("name")
	if name == "" {
		control.FailByMsg("参数名称不能为空")
		return
	}

	r := config.GetParam(name, "")
	control.SuccessData(r)
}

// ModifyParam doc
//
// @tags        服务参数管理
// @Summary     修改参数内容
// @Description 修改参数内容
// @Produce     json
// @Param       data body     model.ModifyParamInfo true "修改参数的内容"
// @Success     200  {object} res.Base              "code: 200 成功"
// @Failure     500  {object} res.Base              "错误返回内容"
// @Router      /granada/api/v1/param [put]
// @Security    ApiKeyAuth
func (ParamsController) ModifyParam(ctx *gin.Context) {
	control := New(ctx)

	info := new(config.Param)
	err := control.Bind(info)
	if err != nil {
		control.FailByMsgf("绑定参数失败，失败原因：%s", err.Error())
		return
	}

	r := config.SetParam(info.Name, info.Value, info.Mark)
	if err != nil {
		control.FailByMsgf("修改配置参数失败，失败原因：%s", err.Error())
		return
	}

	// 设置GOMAXPROCS
	if strings.ToUpper(info.Name) == "GOMAXPROCS" {
		runtime.GOMAXPROCS(r.Int())
	}

	control.Success()
}
