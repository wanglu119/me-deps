package main

import (
	"time"
	
	depDbUtils "github.com/wanglu119/me-deps/db/common/utils"
)

type Test struct {
	Id uint64 			`xorm:"pk autoincr" json:"id"`
	Username string 	`bson:"username,omitempty" xorm:"index" json:"username"`
	CreateAt depDbUtils.Time	`bson:"create_at,omitempty" json:"createAt"`
	
	// front param
	Distance float32	`xorm:"-" json:"distance"`
}

func (t *Test) BeforeInsert() {
	t.CreateAt = depDbUtils.Time(time.Now())
}
