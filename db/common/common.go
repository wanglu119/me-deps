package common

import (

)

type PageInfo struct {
	Limit int
	Skip int
}

type SortInfo struct {
	FieldName string 
	SortType string
}

const (
	SortDesc = "des"
	SortAsc = "aes"
)

type CommonOp interface {
	InsertOne(bean interface{}) error
	FindOne(bean interface{}) error
	Find(beans interface{},condiBean interface{}) error
	FindByPage(beans interface{},condiBean interface{}, page *PageInfo) (int64,error)
	UpdateOne(bean interface{}, condiBean interface{}) error
	DeleteOne(bean interface{}) error 
	DeleteMany(bean interface{}) error
	Count(bean interface{}) (int64, error)
	FindWithSort(beans interface{},condiBean interface{}, sort *SortInfo) error
	FindByPageWithSort(beans interface{},condiBean interface{}, page *PageInfo, sort *SortInfo) (int64,error)
}

