package controller

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/issueye/log-agent/internal/global"
	"github.com/issueye/log-agent/internal/model"
	"github.com/issueye/log-agent/internal/service"
	"github.com/jinzhu/copier"
)

type TBZDDSRRController struct {
	Controller
}

func NewTBZDDSRRController() *TBZDDSRRController {
	return new(TBZDDSRRController)
}

// List doc
//
// @tags        定时任务管理
// @Summary     获取定时任务列表
// @Description 获取定时任务列表
// @Produce     json
// @Param       isNotPaging query    string                           false "是否需要分页， 默认需要， 如果不分页 传 true"
// @Param       pageNum     query    string                           true  "页码， 如果不分页 传 0"
// @Param       pageSize    query    string                           true  "一页大小， 如果不分页 传 0"
// @Param       crwmc       query    string                           false "任务名称"
// @Param       desc        query    string                           false "任务描述"
// @Success     200         {object} res.Full{data=[]models.TBZDDSRR} true  "code: 200 成功"
// @Failure     500         {object} res.Base                         true  "错误返回内容"
// @Router      /granada/api/v1/job/list [get]
// @Security    ApiKeyAuth
func (TBZDDSRRController) List(ctx *gin.Context) {
	control := New(ctx)

	req := new(model.QueryTask)
	err := control.Bind(req)
	if err != nil {
		global.Log.Errorf("绑定请求内容失败，失败原因：%s", err.Error())
		control.FailBind(err)
		return
	}

	list, err := service.NewTask(global.DB).List(req)
	if err != nil {
		global.Log.Errorf("查询定时任务列表失败，失败原因： %s", err.Error())
		control.FailByMsg("查询定时任务列表失败")
		return
	}

	control.SuccessAutoData(req, list)
}

// Create doc
//
// @tags        定时任务管理
// @Summary     添加定时任务数据
// @Description 添加定时任务数据
// @Produce     json
// @Param       data body     model.CreateTBZDDSRR true "添加分诊台字典数据"
// @Success     200  {object} res.Base             true "code: 200 成功"
// @Failure     500  {object} res.Base             true "错误返回内容"
// @Router      /granada/api/v1/job/create [post]
// @Security    ApiKeyAuth
func (TBZDDSRRController) Create(ctx *gin.Context) {
	control := New(ctx)

	req := new(model.CreateTask)
	err := control.Bind(req)
	if err != nil {
		global.Log.Errorf("绑定参数失败，失败原因：%s", err.Error())
		control.FailBind(err)
		return
	}

	err = service.NewTask(global.DB).Create(req)
	if err != nil {
		global.Log.Errorf("添加定时任务失败，失败原因：%s", err.Error())
		control.FailByMsgf("添加定时任务失败，失败原因：%s", err.Error())
		return
	}
	control.Success()
}

// Modify doc
//
// @tags        定时任务管理
// @Summary     修改定时任务数据
// @Description 修改定时任务数据
// @Produce     json
// @Param       id   path     string               true "id"
// @Param       data body     model.ModifyTBZDDSRR true "添加分诊台字典数据"
// @Success     200  {object} res.Base             true "code: 200 成功"
// @Failure     500  {object} res.Base             true "错误返回内容"
// @Router      /granada/api/v1/job/modify/{id} [put]
// @Security    ApiKeyAuth
func (TBZDDSRRController) Modify(ctx *gin.Context) {
	control := New(ctx)

	req := new(model.ModifyTask)
	err := ctx.Bind(req)
	if err != nil {
		global.Log.Errorf("绑定参数失败，失败原因：%s", err.Error())
		control.FailBind(err)
		return
	}

	id := control.Param("id")
	if id == "" {
		global.Log.Error("修改定时任务ID不能为空")
		control.FailByMsg("修改定时任务ID不能为空")
		return
	}

	// 查询定时任务
	info, err := service.NewTask(global.DB).GetById(id)
	if err != nil {
		global.Log.Errorf("查询定时任务失败，失败原因：%s", err.Error())
		control.FailByMsg("查询定时任务失败")
		return
	}

	// 系统任务不允许修改表述信息
	if info.Mark == global.SYS_AUTO_CREATE {
		// 判断描述是否被修改
		if !strings.EqualFold(info.Mark, req.Mark) {
			global.Log.Errorf("定时任务【%s-%s】由系统生成，不允许修改描述信息", info.Name, info.ID)
			control.FailByMsgf("定时任务【%s-%s】由系统生成，不允许修改描述信息", info.Name, info.ID)
			return
		}
	}

	err = service.NewTask(global.DB).Modify(id, req)
	if err != nil {
		global.Log.Errorf("修改定时任务信息失败，失败原因：%s", err.Error())
		control.FailByMsg("修改定时任务信息失败")
		return
	}

	// 查询定时任务
	info, err = service.NewTask(global.DB).GetById(id)
	if err != nil {
		global.Log.Errorf("查询定时任务失败，失败原因：%s", err.Error())
		control.FailByMsg("查询定时任务失败")
		return
	}

	notice := new(model.NoticeJob)
	err = copier.Copy(notice, info)
	if err != nil {
		global.Log.Errorf("拷贝定时任务信息失败，失败原因：%s", err.Error())
		control.FailByMsg("拷贝定时任务信息失败")
		return
	}

	if info.State {
		notice.Type = model.MODIFY_JOB
		global.JobChan <- notice
	}

	control.Success()
}

// ModifyStatus doc
//
// @tags        定时任务管理
// @Summary     修改定时任务数据
// @Description 修改定时任务数据
// @Produce     json
// @Param       id  path     string   true "id"
// @Success     200 {object} res.Base true "code: 200 成功"
// @Failure     500 {object} res.Base true "错误返回内容"
// @Router      /granada/api/v1/job/modifyStatus/{id} [put]
// @Security    ApiKeyAuth
func (TBZDDSRRController) ModifyStatus(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		global.Log.Errorf("修改定时任务状态，参数ID不能为空")
		control.FailByMsg("修改定时任务状态，参数ID不能为空")
		return
	}

	// 获取当前定时任务的状态
	job, err := service.NewTask(global.DB).GetById(id)
	if err != nil {
		global.Log.Errorf("查询定时任务信息失败，失败原因：%s", err.Error())
		control.FailByMsgf("查询定时任务信息失败，失败原因：%s", err.Error())
		return
	}

	err = service.NewTask(global.DB).ModifyStatus(id, !job.State)
	if err != nil {
		global.Log.Errorf("修改定时任务信息失败，失败原因：%s", err.Error())
		control.FailByMsgf("修改定时任务信息失败，失败原因：%s", err.Error())
		return
	}

	// 获取当前定时任务的状态
	job, err = service.NewTask(global.DB).GetById(id)
	if err != nil {
		global.Log.Errorf("查询定时任务信息失败，失败原因：%s", err.Error())
		control.FailByMsgf("查询定时任务信息失败，失败原因：%s", err.Error())
		return
	}

	// 传递到管道
	notice := new(model.NoticeJob)
	err = copier.Copy(notice, job)
	if err != nil {
		global.Log.Errorf("拷贝定时任务信息失败，失败原因：%s", err.Error())
		control.FailByMsgf("拷贝定时任务信息失败，失败原因：%s", err.Error())
		return
	}

	// 开启定时任务则是添加定时任务，否则就是删除定时任务
	if job.State {
		notice.Type = model.ADD_JOB
	} else {
		notice.Type = model.DEL_JOB
	}

	global.JobChan <- notice

	control.Success()
}

// BatchDelete doc
//
// @tags        定时任务管理
// @Summary     批量删除定时任务数据
// @Description 批量删除定时任务数据
// @Produce     json
// @Param       ids body     model.DeleteTBZDDSRR true "ids"
// @Success     200 {object} res.Base             true "code: 200 成功"
// @Failure     500 {object} res.Base             true "错误返回内容"
// @Router      /granada/api/v1/job/batchDelete [delete]
// @Security    ApiKeyAuth
func (TBZDDSRRController) BatchDelete(ctx *gin.Context) {
	control := New(ctx)

	req := new(model.DelTask)
	err := ctx.Bind(req)
	if err != nil {
		global.Log.Errorf("绑定请求参数失败，失败原因：%s", err.Error())
		control.FailBind(err)
		return
	}

	err = service.NewTask(global.DB).BatchDelete(req)
	if err != nil {
		global.Log.Errorf("删除定时任务失败，失败原因：%s", err.Error())
		control.FailByMsgf("删除定时任务失败，失败原因：%s", err.Error())
		return
	}

	control.Success()
}

// Delete doc
//
// @tags        定时任务管理
// @Summary     删除定时任务数据
// @Description 删除定时任务数据
// @Produce     json
// @Param       id  path     string   true "id"
// @Success     200 {object} res.Base true "code: 200 成功"
// @Failure     500 {object} res.Base true "错误返回内容"
// @Router      /granada/api/v1/job/delete/{id} [delete]
// @Security    ApiKeyAuth
func (TBZDDSRRController) Delete(ctx *gin.Context) {
	control := New(ctx)

	id := control.Param("id")
	if id == "" {
		global.Log.Error("删除定时任务，编码不能为空")
		control.FailBind(errors.New("[id]不能为空"))
		return
	}

	// 查询定时任务
	info, err := service.NewTask(global.DB).GetById(id)
	if err != nil {
		global.Log.Errorf("查询定时任务失败，失败原因：%s", err.Error())
		control.FailByMsgf("查询定时任务失败，失败原因：%s", err.Error())
		return
	}

	if info.Mark == global.SYS_AUTO_CREATE {
		global.Log.Errorf("定时任务【%s-%s】由系统生成，不允许删除", info.Name, info.ID)
		control.FailByMsgf("定时任务【%s-%s】由系统生成，不允许删除", info.Name, info.ID)
		return
	}

	err = service.NewTask(global.DB).Delete(id)
	if err != nil {
		global.Log.Errorf("删除定时任务失败，失败原因：%s", err.Error())
		control.FailByMsgf("删除定时任务失败，失败原因：%s", err.Error())
		return
	}

	// 传递到管道
	notice := new(model.NoticeJob)
	err = copier.Copy(notice, info)
	if err != nil {
		global.Log.Errorf("拷贝定时任务信息失败，失败原因：%s", err.Error())
		control.FailByMsgf("拷贝定时任务信息失败，失败原因：%s", err.Error())
		return
	}

	// 如果定时任务是在运行状态则需要删除定时任务
	if info.State {
		notice.Type = model.DEL_JOB
	}

	global.JobChan <- notice

	control.Success()
}
