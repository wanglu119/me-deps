package mgo

import (
	"reflect"
)

type Test struct {
	Common `bson:"-"`
}

func (t *Test) CreateIndex() {
	CreateUniqueIndex(reflect.ValueOf(t).Elem().Type(), "name")
}

/*
func (t *Test) GetAll() []*common.Test {
	beans := []*common.Test{}
	t.find(&beans, bson.D{{}})
	
	return beans; 
}
*/

