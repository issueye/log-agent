package service

import (
	"fmt"
	"strconv"

	"github.com/issueye/log-agent/internal/model"
	"github.com/issueye/log-agent/pkg/utils"
	"gorm.io/gorm"
)

type Monitor struct {
	Db *gorm.DB
	*BaseService
}

func NewMonitor(db *gorm.DB) *Monitor {
	monitor := new(Monitor)
	monitor.Db = db
	monitor.BaseService = NewBaseService(db)
	return monitor
}

func (srv *Monitor) Query(req *model.QueryMonitor) ([]*model.Monitor, error) {
	list := make([]*model.Monitor, 0)
	err := srv.DataFilter(model.Monitor{}.TableName(), req, &list, func(db *gorm.DB) (*gorm.DB, error) {
		query := db.Order("id")

		if req.Condition != "" {
			query = query.Where("name like ?", fmt.Sprintf("%%%s%%", req.Condition))
			query = query.Or("log_path like ?", fmt.Sprintf("%%%s%%", req.Condition))
			query = query.Or("script_path like ?", fmt.Sprintf("%%%s%%", req.Condition))
		}

		return query, nil
	})
	return list, err
}

func (srv *Monitor) GetById(id string) (*model.Monitor, error) {
	data := new(model.Monitor)
	err := srv.Db.Model(data).Where("id = ?", id).Find(data).Error
	return data, err
}

// 创建数据
func (srv *Monitor) Create(data *model.CreateMonitor) error {
	m := new(model.Monitor)
	m.ID = strconv.FormatInt(utils.GenID(), 10)
	m.Name = data.Name
	m.LogPath = data.LogPath
	m.ScriptPath = data.ScriptPath
	m.Level = data.Level
	m.State = false // 初始添加时，状态为未开启
	m.CreateTime = utils.GetNowStr()
	return srv.Db.Create(m).Error
}

func (srv *Monitor) Modify(data *model.ModifyMonitor) error {
	m := make(map[string]any)
	m["name"] = data.Name
	m["log_path"] = data.LogPath
	m["level"] = data.Level
	m["script_path"] = data.ScriptPath

	return srv.Db.Model(&model.Monitor{}).Updates(m).Error
}

func (srv *Monitor) ModifyState(id string) (bool, error) {
	data := new(model.Monitor)
	err := srv.Db.Model(&model.Monitor{}).Where("id = ?", id).Find(data).Error
	if err != nil {
		return false, err
	}

	nowState := !data.State
	err = srv.Db.Model(data).Update("state = ?", nowState).Error
	if err != nil {
		return false, err
	}

	return nowState, nil
}

func (srv *Monitor) DelMonitor(id string) error {
	return srv.Db.Where("id = ?", id).Delete(&model.Monitor{}).Error
}
