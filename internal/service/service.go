package service

import (
	"errors"
	"reflect"

	"gorm.io/gorm"
)

type BaseService struct {
	Db *gorm.DB
}

func NewBaseService(db *gorm.DB) *BaseService {
	return &BaseService{
		Db: db,
	}
}

type ListFilter func(db *gorm.DB) (*gorm.DB, error)

// DataFilter 数据过滤
func (srv BaseService) DataFilter(tableName string, req, list interface{}, f ListFilter) error {
	query := srv.Db.Table(tableName)
	db, err := f(query)
	if err != nil {
		return err
	}
	count := int64(0)
	err = db.Count(&count).Error
	if err != nil {
		return err
	}
	ref := reflect.ValueOf(req)
	// 判断 ref 是否是 ptr 类型
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}

	// 判断 req 是否有 Total 字段
	// 如果有则将count 赋值给 req.Total
	if !ref.FieldByName("Total").IsValid() {
		return errors.New("请求参数错误，传入参数中没有 Total 字段")
	}
	ref.FieldByName("Total").SetInt(count)

	// 判断是否需要进行分页
	pageNum := ref.FieldByName("PageNum").Int()
	pageSize := ref.FieldByName("PageSize").Int()
	// 如果 页码、每页大小小于0 则也不分页
	if pageNum == 0 || pageSize == 0 {
		err = db.Find(list).Error
		if err != nil {
			return err
		}
		return nil
	} else {
		// 从 ref 中获取 分页参数 PageNum PageSize
		err = db.Offset(int(pageNum-1) * int(pageSize)).Limit(int(pageSize)).Find(list).Error
		if err != nil {
			return err
		}
		return nil
	}
}
