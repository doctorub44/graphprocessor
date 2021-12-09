package graphproc

import (
	"reflect"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger
var fastlogger *zap.Logger

//FunctionName :
func FunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

//BaseFunctionName :
func BaseFunctionName(i interface{}) string {
	return strings.Split(FunctionName(i), ".")[1]
}

//CallerName :
func CallerName() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(2, fpcs)
	if n == 0 {
		return "<nil>"
	}
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "<nil>"
	}
	return fun.Name()
}

//Logger : simple structured logger
func Logger(s string, arg ...interface{}) {
	if logger == nil {
		initLogger()
	}
	logger.Info(s+" ", arg)
}

//FastLogger : fast structured logger for hot paths
func FastLogger(s string, arg ...interface{}) {
	if fastlogger == nil {
		initLogger()
	}

	if arg != nil {
		fastlogger.Info(s, zap.Time("timestamp", arg[0].(time.Time)))
	} else {
		fastlogger.Info(s)
	}
}

//initLogger : initialize the logging functions
func initLogger() {
	var err error
	output := []string{"logger.txt"}
	config := zap.NewDevelopmentConfig()
	config.DisableCaller = true
	config.OutputPaths = output
	fastlogger, err = config.Build()
	if err != nil {
		panic(err)
	}
	logger = fastlogger.Sugar()
}
