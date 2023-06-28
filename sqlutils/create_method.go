package sqlutils

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/copier"
	"regexp"
)

func (in InitData) CreateObject() (err error) {
	if in.DB == nil {
		return errors.New("DB 参数不能为空")
	}
	if in.ReqObj == nil {
		return errors.New("RespObj 参数不能为空")
	}
	if in.RespObj == nil {
		return errors.New("RespObj 参数不能为空")
	}

	if err = copier.Copy(in.RespObj, in.ReqObj); err != nil {
		return err
	}

	err = in.DB.Create(in.RespObj).Error
	if err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if ok == true && mysqlError.Number == 1062 {
			compileRegex := regexp.MustCompile("Duplicate entry '(.*-)?(.*?)' for key .*")
			matchArrStr := compileRegex.FindStringSubmatch(err.(*mysql.MySQLError).Message)
			return errors.New(fmt.Sprintf("%s 已存在", matchArrStr[len(matchArrStr)-1]))
		}
		return err
	}

	return nil
}
