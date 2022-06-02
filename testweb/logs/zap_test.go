package logs

import "testing"

func TestSimpleHttpGet(t *testing.T) {
	InitLogger()
	defer sugarLogger.Sync()
	SimpleHttpGet("www.sogo.com")
	SimpleHttpGet("http://www.sogo.com")
	Debug("123")
	Error("123")
	Info("123")
}
