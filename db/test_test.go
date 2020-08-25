package db

import (
	"testing"
	"time"
	
	"github.com/wanglu119/me-deps/db/common"
	"github.com/wanglu119/me-deps/db/common/config"
	"github.com/wanglu119/me-deps/db/mgo"
	"github.com/wanglu119/me-deps/db/orm"
)

func init() {
	test := "mongo"
//	test := "sqlite3"

	if test == "mongo" {
		config.Database.Name = "wl_test"
//		config.Database.Host = "192.168.0.106:27017"
		config.Database.Host = "172.17.0.1:27017"
		config.Database.Type = "mongo"
		
		mgo.AddTable(new(mgo.Test))
		
		err := mgo.NewEngine()
		if err != nil {
			panic(err)
		}
		
		Tests = &tests {
			CommonOp: &mgo.Common{},
		}
	} 
	if test == "sqlite3" {
		config.Database.Type = "sqlite3"
		config.Database.Path = "./bin/test.db"
		
		orm.AddTable(new(common.Test))
		err := orm.NewEngine()
		if err != nil {
			panic(err)
		}
		
		Tests = &tests {
			CommonOp: &orm.Common{},
		}
	}
}


func TestDbMethod(t *testing.T) {
	// insert
	opts := CreateTestOpts{
		Name: "test-"+time.Now().Format("2006-01-02 15:04:05"),
	}
	_, err := Tests.Create(opts)
	if err != nil {
		t.Error(err)
		return
	}
	
	// get all
	ts,err := Tests.GetAll()
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("get all:", len(ts))
	for _, te := range ts {
		t.Log(te.Name)
	}
	
	name := ts[0].Name
	
	// find with condiBean
	condiBean := &common.Test{
		Name: name,
	}
	ts, err = Tests.GetWithCondiBean(condiBean)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("get with condiBean:", len(ts))
	for _, te := range ts {
		t.Log(te.Name)
	}
	
	// find with sort
	sort := &common.SortInfo{
		FieldName: "create_at",
		SortType: common.SortDesc,
	}
	ts, err = Tests.GetWithCondiBeanAndSort(nil, sort)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("get with sort desc:", len(ts))
	for _, te := range ts {
		t.Log(te.Name)
	}
	
	sort = &common.SortInfo{
		FieldName: "create_at",
		SortType: common.SortAsc,
	}
	ts, err = Tests.GetWithCondiBeanAndSort(nil, sort)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("get with sort asc:", len(ts))
	for _, te := range ts {
		t.Log(te.Name)
	}
	
	// page find
	page := &common.PageInfo{
		Limit: 2,
		Skip: 0,
	}
	ts,_, err = Tests.GetByPage(&common.Test{}, page)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("get with page:", len(ts))
	for _, te := range ts {
		t.Log(te.Name)
	}
	
	// page find with sort
	sort = &common.SortInfo{
		FieldName: "create_at",
		SortType: common.SortDesc,
	}
	ts,count, err := Tests.GetByPageWithSort(&common.Test{}, page,sort)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("get page with sort desc:", len(ts), count)
	for _, te := range ts {
		t.Log(te.Name, te.CreateAt.Format("2006-01-02 15:04:05"))
	}
	
	sort = &common.SortInfo{
		FieldName: "create_at",
		SortType: common.SortAsc,
	}
	ts,count, err = Tests.GetByPageWithSort(&common.Test{}, page,sort)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("get page with sort asc:", len(ts), count)
	for _, te := range ts {
		t.Log(te.Name, te.CreateAt.Format("2006-01-02 15:04:05"))
	}
	
	// find one 
	getT,err := Tests.GetByName(name)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("find test: %s, %s", getT.Name, getT.CreateAt.Format("2006-01-02 15:04:05"))
	
	// update one
	getT.CreateAt = time.Now()
	err = Tests.Update(getT, condiBean)
	if err != nil {
		t.Error(err)
		return
	}
	
	// page find with condiBean
	ts,_, err = Tests.GetByPage(condiBean, page)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("get with page condiBean:", len(ts))
	for _, te := range ts {
		t.Log(te.Name,te.CreateAt.Format("2006-01-02 15:04:05"))
	}
	
	// delete
	delBean := &common.Test{
		Name: getT.Name,
	}
	err = Tests.Delete(delBean)
	if err != nil {
		t.Log(err)
		return
	}
}


