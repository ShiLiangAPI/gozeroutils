package sqlutils

import (
	"errors"
	"fmt"
	"github.com/ShiLiangAPI/goutils/function"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/copier"
	"reflect"
	"regexp"
)

// UpdateObject 通过ID修改对象
func (in InitData) UpdateObject() (err error) {
	if err = in.GetObject(); err != nil {
		return err
	}
	if err = copier.Copy(in.RespObj, in.ReqObj); err != nil {
		return err
	}

	if err = in.DB.Updates(in.RespObj).Error; err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if ok && mysqlError.Number == 1062 {
			compileRegex := regexp.MustCompile("Duplicate entry '(.*-)?(.*?)' for key .*")
			matchArrStr := compileRegex.FindStringSubmatch(err.(*mysql.MySQLError).Message)
			return errors.New(fmt.Sprintf("%s 已存在", matchArrStr[len(matchArrStr)-1]))
		}
		return err
	}

	return nil
}

// UpdateFilterObject 通过条件(修改第一个或批量修改)对象
func (in InitData) UpdateFilterObject() (err error) {
	if in.DB == nil {
		return errors.New("DB 参数不能为空")
	}
	if in.ReqObj == nil {
		return errors.New("DB 参数不能为空")
	}

	if in.RespObj != nil {
		if err = in.GetFirstObject(); err != nil {
			return err
		}
		if err := copier.Copy(in.RespObj, in.ReqObj); err != nil {
			return err
		}
		if err = in.query.Updates(in.RespObj).Error; err != nil {
			mysqlError, ok := err.(*mysql.MySQLError)
			if ok == true && mysqlError.Number == 1062 {
				compileRegex := regexp.MustCompile("Duplicate entry '(.*-)?(.*?)' for key .*")
				matchArrStr := compileRegex.FindStringSubmatch(err.(*mysql.MySQLError).Message)
				return errors.New(fmt.Sprintf("%s 已存在", matchArrStr[len(matchArrStr)-1]))
			}
			return err
		}
	} else if in.RespObjList != nil {
		if _, err = in.GetAllObject(); err != nil {
			return err
		}
		for _, respObj := range in.RespObjList.([]any) {
			if err := copier.Copy(respObj, in.ReqObj); err != nil {
				return err
			}
		}
		if err = in.query.Updates(in.RespObjList).Error; err != nil {
			mysqlError, ok := err.(*mysql.MySQLError)
			if ok == true && mysqlError.Number == 1062 {
				compileRegex := regexp.MustCompile("Duplicate entry '(.*-)?(.*?)' for key .*")
				matchArrStr := compileRegex.FindStringSubmatch(err.(*mysql.MySQLError).Message)
				return errors.New(fmt.Sprintf("%s 已存在", matchArrStr[len(matchArrStr)-1]))
			}
			return err
		}
	} else {
		return errors.New("RespObj 或 RespObjList 参数不能为空")
	}

	return nil
}

// UpdateRelationObject T 关联的类型
// mapValue = {"id/pk": "post_id", "update_id": "user_id", "filter_field": "PostID", "update_field": "UserID"}
func (in InitData) UpdateRelationObject() (err error) {
	if in.PK == 0 {
		return errors.New("PK 参数不能为空")
	}
	if in.PKList == nil {
		return errors.New("PKList 参数不能为空")
	}
	if in.FilterField == "" {
		return errors.New("FilterField 参数不能为空")
	}
	if in.UpdateField == "" {
		return errors.New("UpdateField 参数不能为空")
	}
	var modelType any
	if in.RespObj != nil {
		modelType = in.RespObj
	} else if in.Model != nil {
		modelType = in.Model
	} else {
		return errors.New("参数 RespObj 或 Model 不能为空")
	}

	in.query = in.getQuerySet()
	var activeIDList []int64
	var addingIDList []int64
	var deletingIDList []int64
	if err := in.query.Where("id = ?", in.PK).Pluck(in.UpdateField, &activeIDList).Error; err != nil {
		return err
	}
	addingIDList = function.Different[int64](in.PKList, activeIDList)
	deletingIDList = function.Different[int64](activeIDList, in.PKList)
	if len(deletingIDList) > 0 {
		if err = in.query.Unscoped().Where(deletingIDList).Delete(modelType).Error; err != nil {
			return err
		}
	}
	if len(addingIDList) > 0 {
		var relationList []any

		for _, value := range addingIDList {
			modelReflect := reflect.ValueOf(modelType).Elem()
			modelReflect.FieldByName(in.FilterField).SetInt(in.PK)
			modelReflect.FieldByName(in.UpdateField).SetInt(value)
			relationList = append(relationList, modelType)
		}
		if err = in.query.CreateInBatches(relationList, len(relationList)).Error; err != nil {
			return err
		}
	}

	return nil
}
