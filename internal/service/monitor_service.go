package service

import (
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
