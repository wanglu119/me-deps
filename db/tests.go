package db

import (
	"github.com/wanglu119/me-deps/db/common"
	"github.com/wanglu119/me-deps/db/common/config"
	"github.com/wanglu119/me-deps/db/mgo"
	"github.com/wanglu119/me-deps/db/orm"
)

func NewTest() *common.Test {
	var op common.TestOp
	if config.Database.Type == "mongo" {
		op = &mgo.Test{}
	} else {
		op = &orm.Test{}
	}
	
	return &common.Test{TestOp: op}
}

type CreateTestOpts struct {
	Name string
}

type TestsStore interface {
	Create(opts CreateTestOpts) (*common.Test, error)
	GetByName(name string) (*common.Test, error)
	GetAll() ([]*common.Test, error)
	GetWithCondiBean(condiBean *common.Test) ([]*common.Test, error)
	GetByPage(condiBean *common.Test, page *common.PageInfo) ([]*common.Test,int64, error)
	Update(bean *common.Test, condiBean *common.Test) error
	Delete(bean *common.Test) error
	GetWithCondiBeanAndSort(condiBean *common.Test, sort *common.SortInfo) ([]*common.Test, error)
	GetByPageWithSort(condiBean *common.Test, page *common.PageInfo, sort *common.SortInfo) ([]*common.Test, int64, error)
}

var Tests TestsStore

type tests struct {
	common.CommonOp
}

func (ts *tests) setTestOp(t *common.Test) {
	if config.Database.Type == "mongo" {
		t.TestOp = &mgo.Test{}
	} else {
		t.TestOp = &orm.Test{}
	}
}

func (ts *tests) Create(opts CreateTestOpts) (*common.Test, error) {
	
	t := &common.Test {
		Name: opts.Name,
	}
	
	err := ts.InsertOne(t)
	
	if err != nil {
		return nil, err
	}
	
	ts.setTestOp(t)
	
	return t, nil
}

func (ts *tests) GetByName(name string) (*common.Test, error) {
	t := &common.Test {
		Name: name,
	}
	err := ts.FindOne(t)
	if err != nil {
		return nil, err
	}
	
	ts.setTestOp(t)
	
	return t, nil
}

func (ts *tests) GetAll() ([]*common.Test, error) {
	beans := []*common.Test{}
	err := ts.Find(&beans, nil)
	if err != nil {
		return nil, err
	}
	return beans, nil
}

func (ts *tests) GetWithCondiBean(condiBean *common.Test) ([]*common.Test, error) {
	beans := []*common.Test{}
	err := ts.Find(&beans, condiBean)
	if err != nil {
		return nil, err
	}
	return beans, nil
}

func (ts *tests) GetByPage(condiBean *common.Test, page *common.PageInfo) ([]*common.Test, int64, error) {
	beans := []*common.Test{}
	count, err := ts.FindByPage(&beans, condiBean, page)
	if err != nil {
		return nil,0, err
	}
	return beans,count, nil
}

func (ts *tests) Update(bean *common.Test, condiBean *common.Test) error {
	return ts.UpdateOne(bean, condiBean)
}

func (ts *tests) Delete(bean *common.Test) error {
	return ts.DeleteOne(bean)
}

func (ts *tests) GetWithCondiBeanAndSort(condiBean *common.Test, sort *common.SortInfo) ([]*common.Test, error) {
	beans := []*common.Test{}
	err := ts.FindWithSort(&beans, condiBean, sort)
	if err != nil {
		return nil, err
	}
	return beans, nil
}

func (ts *tests) GetByPageWithSort(condiBean *common.Test, page *common.PageInfo, sort *common.SortInfo) ([]*common.Test, int64, error) {
	beans := []*common.Test{}
	count, err := ts.FindByPageWithSort(&beans, condiBean, page, sort)
	if err != nil {
		return nil,0, err
	}
	return beans,count, nil
}
