package main

import (
	wlog "github.com/wanglu119/me-deps/log"
)

func main() {
	log := wlog.GetLogger()
	log.Info("test info log")
	log.Warn("test warn log")
	log.Error("test error log")
}
