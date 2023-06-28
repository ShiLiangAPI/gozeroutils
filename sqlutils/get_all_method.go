package sqlutils

import (
	"github.com/pkg/errors"
)

// GetPageObject 获取分页数据
//
//	"where": map[string]any{
//		 "name Like ?": req.Search,
//	},
func (in InitData) GetPageObject() (resp map[string]any, err error) {
	if in.DB == nil {
		return nil, errors.New("DB 参数不能为空")
	}
	if in.RespObjList == nil {
		return nil, errors.New("RespObj 参数不能为空")
	}
	in.query = in.getQuerySet()
	in.query = in.filterData()
	// 处理分页
	pageMap := in.paginate()
	// 查询结果
	if err = in.query.Find(in.RespObjList).Error; err != nil {
		return nil, err
	}
	pageMap["list"] = in.RespObjList
	return pageMap, nil
}

// GetAllObject 获取不分页所有数据
func (in InitData) GetAllObject() (resp map[string]any, err error) {
	if in.DB == nil {
		return nil, errors.New("DB 参数不能为空")
	}
	if in.RespObjList == nil {
		return nil, errors.New("RespObj 参数不能为空")
	}
	in.query = in.getQuerySet()
	in.query = in.filterData()
	// 查询结果
	if err = in.query.Find(in.RespObjList).Error; err != nil {
		return nil, err
	}
	return map[string]any{"list": in.RespObjList}, nil
}

// GetTreeObject
//
//	"where": map[string]any{
//		 "name Like ?": req.Search,
//	},
//func (in InitData) GetTreeObject() (resp map[string]any, err error) {
//
//	if in.DB == nil {
//		return nil, errors.New("DB 参数不能为空")
//	}
//	if in.RespObjList == nil {
//		return nil, errors.New("RespObj 参数不能为空")
//	}
//	in.query = in.getQuerySet()
//	in.query = in.filterData()
//
//	getTree := func(obj any) {
//		immutable := reflect.ValueOf(obj)
//		id := immutable.FieldByName("ID")
//		in.query.Where("parent_id = ?", id)
//	}
//
//	if in.NodeId == 0 {
//		_, err := in.GetAllObject()
//		if err != nil {
//			return nil, err
//		}
//	} else {
//		if err = in.query.Where("id = ?", in.NodeId).Find(in.RespObjList).Error; err != nil {
//			return nil, err
//		}
//	}
//
//	for _, value := range in.RespObjList {
//		getTree(value)
//	}
//
//	return map[string]any{"list": in.RespObjList}, nil
//}

// GetRelationObject
// mapValue=mapValue = {"id/pk": "post_id", "filter_id": "user_id", "where": {}， "filter_where": {}}
func (in InitData) GetRelationObject() error {
	if in.PK == 0 {
		return errors.New("PK 参数不能为空")
	}
	if in.RelationField == "" {
		return errors.New("RelationField 参数不能为空")
	}
	if in.RespObjList == nil {
		return errors.New("RespObjList 参数不能为空")
	}
	in.query = in.getQuerySet()
	querySet := in.filterData()

	var IDList []int64
	querySet.Where("id = ?", in.PK).Pluck(in.RelationField, &IDList)

	if len(IDList) > 0 {
		in.Where = in.RelationWhere
		filterQuerySet := in.filterData()
		if err := filterQuerySet.Where(IDList).Scan(&in.RespObjList).Error; err != nil {
			return err
		}
	}

	return nil
}

///////////////////////////// 不对外使用 /////////////////////////////

// Paginate 默认分页
func (in InitData) paginate() map[string]any {
	if in.NoPage {
		return map[string]any{}
	}
	if in.CurrentPage <= 0 {
		in.CurrentPage = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	var total int64
	in.query.Count(&total)
	in.query = in.query.Offset((in.CurrentPage - 1) * in.PageSize).Limit(in.PageSize)
	return map[string]any{
		"page_num":  in.CurrentPage,
		"page_size": in.PageSize,
		"total":     total,
	}
}
