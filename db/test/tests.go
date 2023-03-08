package main

import (
	depDbCommon "github.com/wanglu119/me-deps/db/common"
)

type TestStore interface {
	depDbCommon.CommonOp
}

var Tests TestStore

// -------------------------------------------------

type tests struct {
	depDbCommon.CommonOp
}
