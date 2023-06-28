package sqlutils

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

// GetObject 通过id获取对象
func (in InitData) GetObject() (err error) {
	if in.DB == nil {
		return errors.New("DB 参数不能为空")
	}
	if in.PK == "" {
		return errors.New("PK 参数不能为空")
	}
	if in.RespObj == nil {
		return errors.New("RespObj 参数不能为空")
	}
	in.query = in.getQuerySet()
	in.query = in.filterData()

	in.query.First(in.RespObj, in.PK)

	return nil
}

// GetFirstObject 通过筛选获取第一个对象
func (in InitData) GetFirstObject() (err error) {
	if in.DB == nil {
		return errors.New("DB 参数不能为空")
	}
	if in.RespObj == nil {
		return errors.New("RespObj 参数不能为空")
	}
	in.query = in.getQuerySet()
	in.query = in.filterData()

	err = in.query.First(in.RespObj).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("数据不存在")
		} else {
			return err
		}
	}
	return nil
}

///////////////////////////// 不对外使用 /////////////////////////////

// 获取带有表信息的 *gorm.DB
func (in InitData) getQuerySet() *gorm.DB {
	if in.RespObj != nil {
		return in.DB.Model(in.RespObj)
	}
	if in.Model != nil {
		return in.DB.Model(in.Model)
	}
	if in.Table != "" {
		return in.DB.Model(in.Model)
	}

	return in.DB
}

func (in InitData) filterData() *gorm.DB {
	if in.Order != "" {
		in.query = in.query.Order(in.Order)
	}
	if in.OrderList != nil {
		for _, order := range in.OrderList {
			if order != "" {
				in.query = in.query.Order(order)
			}
		}
	}
	if in.Preload != "" {
		in.query = in.query.Preload(in.Preload)
	}
	if in.PreloadList != nil {
		for _, preload := range in.PreloadList {
			if preload != "" {
				in.query = in.query.Preload(preload)
			}
		}
	}
	if in.Where != nil {
		in.query = in.queryWhere()
	}

	return in.query
}

func (in InitData) queryWhere() *gorm.DB {
	for condition, args := range in.Where {
		queryValT := reflect.TypeOf(args).Kind()
		if queryValT == reflect.Slice {
			in.query = in.query.Where(handlerConditionString(condition), args.([]any)...)
		} else {
			in.query = in.query.Where(handlerConditionString(condition), args)
		}
	}
	return in.query
}

// handlerConditionString 数据库字段匹配条件；无？自动匹配 = ？
func handlerConditionString(condition string) string {
	var stringBuild strings.Builder
	stringBuild.WriteString(condition)
	if strings.Index(condition, "?") == -1 {
		stringBuild.WriteString(" = ?")
	}
	return stringBuild.String()
}
