package service

import (
	"fmt"
	"strings"

	"github.com/issueye/log-agent/internal/config"
	"github.com/issueye/log-agent/internal/model"
	"github.com/issueye/log-agent/pkg/utils"
	"gorm.io/gorm"
)

type User struct {
	*BaseService
}

func NewUser(db *gorm.DB) *User {
	return &User{
		BaseService: NewBaseService(db),
	}
}

// FindUser
// 查找用户
func (user *User) FindUser(info *model.Login) (*model.User, error) {
	query := user.Db.Model(&model.User{}).Order("id")
	query = query.Where("account = ?", info.Account)

	// 判断是否需要验证密码
	r := config.GetParam("SERVER-MODE", "release")
	fmt.Printf("当前运行模式：%s", r.String())
	if strings.EqualFold(r.String(), "release") {
		query = query.Where("password = ?", info.Password)
	}

	data := new(model.User)
	err := query.Find(data).Error
	return data, err
}

// CreateAdminNonExistent
// 创建管理员用户，如果不存在
func (user *User) CreateAdminNonExistent() error {
	isHave := int64(0)
	err := user.Db.Model(&model.User{}).Where("account = ?", "admin").Count(&isHave).Error
	if err != nil {
		return err
	}

	if isHave == 0 {
		info := new(model.User)
		info.ID = utils.GenID()
		info.Account = "admin"
		info.Name = "管理员"
		info.Password = "123456"
		info.State = 1
		info.Mark = "系统自动生成的管理员数据"
		return user.Db.Create(info).Error
	} else {
		return nil
	}
}

// Create
// 创建用户信息
func (user *User) Create(data *model.CreateUser) error {
	info := new(model.User)
	info.ID = utils.GenID()
	info.Account = data.Account
	info.Name = data.Name
	info.Password = data.Password
	info.Mark = data.Mark
	info.State = 1
	return user.Db.Create(info).Error
}

// GetByAccount
// 查找用户是否存在
func (user *User) GetByAccount(account string) (*model.User, error) {
	info := new(model.User)
	err := user.Db.Model(info).Where("account = ?", account).Find(info).Error
	return info, err
}

// Modify
// 修改用户信息
func (user *User) Modify(info *model.ModifyUser) error {
	m := make(map[string]any)
	m["account"] = info.Account
	m["name"] = info.Name
	m["password"] = info.Password
	m["mark"] = info.Mark

	return user.Db.Model(&model.User{}).Where("id = ?", info.ID).Updates(m).Error
}

// Status
// 修改用户信息
func (user *User) Status(info *model.StatusUser) error {
	return user.Db.
		Model(&model.User{}).
		Where("id = ?", info.ID).
		Update("state", info.State).
		Error
}

// Delete
// 删除用户信息
func (user *User) Delete(id int64) error {
	return user.Db.Where("id = ?", id).Delete(&model.User{}).Error
}

// List
// 获取用户列表
func (user *User) List(info *model.QueryUser) ([]*model.User, error) {
	userInfo := new(model.User)
	list := make([]*model.User, 0)
	err := user.DataFilter(userInfo.TableName(), info, &list, func(db *gorm.DB) (*gorm.DB, error) {
		query := db.Order("id")

		// 用户姓名
		if info.Name != "" {
			query = query.Where("name like ?", fmt.Sprintf("%%%s%%", info.Name))
		}

		return query, nil
	})
	return list, err
}
