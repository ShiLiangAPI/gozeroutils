package sqlutils

import (
	"errors"
	"gorm.io/gorm"
)

// DeleteOneObject T 删除数据的结构体类型
func (in InitData) DeleteOneObject() (err error) {
	if in.DB == nil {
		return errors.New("DB 参数不能为空")
	}
	if in.RespObj == nil {
		return errors.New("RespObj 参数不能为空")
	}
	if in.PK == "" {
		return errors.New("PK 参数不能为空")
	}
	if err = in.GetObject(); err != nil {
		return err
	}
	if err = in.DB.Delete(in.RespObj).Error; err != nil {
		if err == gorm.ErrMissingWhereClause {
			return errors.New("需要删除的数据不存在")
		}
		return err
	}

	return nil
}

func (in InitData) DeleteAllObject() (err error) {
	if in.DB == nil {
		return errors.New("DB 参数不能为空")
	}
	if in.PKList == nil {
		return errors.New("PKList 参数不能为空")
	}
	if in.RespObjList == nil {
		return errors.New("RespObjList 参数不能为空")
	}
	if err = in.DB.Where(in.PKList).Find(in.RespObjList).Error; err != nil {
		return err
	}
	if err = in.DB.Delete(in.RespObjList).Error; err != nil {
		if err == gorm.ErrMissingWhereClause {
			return errors.New("需要删除的数据不存在")
		}
		return err
	}

	return nil
}

func (in InitData) DeleteTreeObject() (err error) {
	if in.DB == nil {
		return errors.New("DB 参数不能为空")
	}
	if in.PKList == nil {
		return errors.New("PKList 参数不能为空")
	}
	in.query = in.getQuerySet()
	var modelType any
	if in.RespObj != nil {
		modelType = in.RespObj
	} else if in.Model != nil {
		modelType = in.Model
	} else {
		return errors.New("参数 RespObj 或 Model 不能为空")
	}

	for {
		if len(in.PKList) <= 0 {
			return nil
		}
		if err = in.query.Where(in.PKList).Delete(modelType).Error; err != nil {
			return err
		}

		parentIDList := in.PKList
		in.PKList = []string{}

		if err = in.query.Where("parent_id IN (?)", parentIDList).Pluck("id", &in.PKList).Error; err != nil {
			return err
		}
	}
}
