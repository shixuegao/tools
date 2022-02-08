package log

import (
	logI "umx/tools/logger"
	zap "umx/tools/logger/condidate/zap"
)

var Logger logI.Logger
var factory *zap.ZapFactory

func init() {
	Logger = zapInit()
}

func zapInit() logI.Logger {
	factory = zap.BetterNewZapFactory("pressure.server.log", logI.Info, 30, 7)
	return factory.Logger()
}

func Clear() {
	factory.Clear()
}
