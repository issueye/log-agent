package service

import (
	"fmt"
	"strconv"

	"github.com/issueye/log-agent/internal/global"
	"github.com/issueye/log-agent/internal/model"
	"github.com/issueye/log-agent/pkg/utils"

	"gorm.io/gorm"
)

type Task struct {
	Db *gorm.DB
	*BaseService
}

func NewTask(db *gorm.DB) *Task {
	task := new(Task)
	task.Db = db
	task.BaseService = NewBaseService(db)
	return task
}

// CreateJobNoExistent
// 初始系统添加的定时任务
func (srv *Task) CreateJobNoExistent() {

}

// CreateJob
// 添加定时任务
func (srv *Task) CreateJob(id, name, path string) {
	table := new(model.Task)
	table.ID = id
	table.Name = name
	table.Path = path
	table.State = false
	table.Mark = global.SYS_AUTO_CREATE
	table.Expression = "0/5 * * * * ?" // 默认为五秒执行一次
	srv.Db.Create(table)
}

// isNotHave
// 是否存在
func (srv *Task) isNotHave(id string) bool {
	count := int64(0)
	err := srv.Db.Model(&model.Task{}).Where("name = ?", id).Count(&count).Error
	if err != nil {
		return true
	}

	return count == 0
}

// Create
// 添加定时任务数据
func (srv *Task) Create(req *model.CreateTask) error {
	table := new(model.Task)
	table.ID = strconv.FormatInt(utils.GenID(), 10)
	table.State = false
	table.Name = req.Name
	table.Path = req.Path
	table.Mark = req.Mark
	table.Expression = req.Expression
	return srv.Db.Create(table).Error
}

func (srv *Task) Delete(id string) error {
	// 删除分诊台字典表
	return srv.Db.Where(fmt.Sprintf("id in (?)", id)).Delete(&model.Task{}).Error
}

func (srv *Task) BatchDelete(ids *model.DelTask) error {
	// 删除分诊台字典表
	return srv.Db.Model(&model.Task{}).Delete("id in (?)", ids.Ids).Error
}

func (srv *Task) Modify(id string, req *model.ModifyTask) error {
	m := make(map[string]interface{})
	m["name"] = req.Name
	m["expression"] = req.Expression
	m["path"] = req.Path
	m["mark"] = req.Mark
	return srv.Db.Model(&model.Task{}).Where("id = ?", id).Updates(&m).Error
}

func (srv *Task) ModifyStatus(id string, state bool) error {
	return srv.Db.Model(&model.Task{}).Where("id = ?", id).Update("state", state).Error
}

func (srv *Task) GetById(id string) (*model.Task, error) {
	data := new(model.Task)
	err := srv.Db.Model(data).Where("id = ?", id).Find(data).Error
	return data, err
}

func (srv *Task) List(req *model.QueryTask) ([]*model.Task, error) {
	list := make([]*model.Task, 0)

	err := srv.DataFilter(model.Task{}.TableName(), req, &list, func(db *gorm.DB) (*gorm.DB, error) {
		query := db.Order("id")

		if req.Name != "" {
			query = query.Where("name like ?", fmt.Sprintf("%%%s%%", req.Name))
		}

		if req.Mark != "" {
			query = query.Where("mark like ?", fmt.Sprintf("%%%s%%", req.Mark))
		}

		return query, nil
	})
	return list, err
}
