package log

import (
	"testing"
)

func TestDefaultUse(t *testing.T) {
	logger.Debug("defaultLoggerUser")
	logger.Info("xxxxxxxxxxxxxxxxxxx")
}

func TestSetFile(t *testing.T) {
	logger.SetFileLocation("/data/test/t3/zap.log")
	logger.Info("ffffffffffffffffffffffffffffffffffffff")
}
