package common

import (
	"time"
	"fmt"
)

type Test struct {
	TestOp `bson:"-" xorm:"-" json:"-"`
	Name string `bson:"name,omitempty" xrom:"pk" json:"name" `
	CreateAt time.Time `bson:"create_at,omitempty" xrom:"create_at", json:"create_at"`
}

func (t *Test) BeforeInsert() {
	t.CreateAt = time.Now()
}

func (t *Test) BeforeUpdate() {
	fmt.Println("-----> BeforeUpdate")
}

type TestOp interface {
}


