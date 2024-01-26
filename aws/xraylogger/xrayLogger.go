package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/cultureamp/ca-go/log"
)

type xrayLogger struct {
	logger *log.Logger
}

func newXrayLogger() *xrayLogger {
	config := log.NewLoggerConfig()
	newLogger := log.NewLogger(config)
	return &xrayLogger{
		logger: newLogger,
	}
}

func (xl xrayLogger) Log(level xraylog.LogLevel, msg fmt.Stringer) {
	props := log.SubDoc().Str("message", msg.String())

	switch level {
	case xraylog.LogLevelDebug:
		xl.logger.Debug("xray_diagnostic").Properties(props).WithSystemTracing().Send()
	case xraylog.LogLevelInfo:
		xl.logger.Info("xray_diagnostic").Properties(props).WithSystemTracing().Send()
	case xraylog.LogLevelWarn:
		xl.logger.Warn("xray_diagnostic").Properties(props).WithSystemTracing().Send()
	case xraylog.LogLevelError:
		xl.logger.Error("xray_diagnostic", errors.New("xray_diagnostic")).Properties(props).WithSystemTracing().Send()
	default:
		xl.logger.Debug("xray_diagnostic").Properties(props).WithSystemTracing().Send()
	}
}

type printArgs struct {
	s string
}

func newPrintArgs(s string) printArgs {
	return printArgs{s: s}
}

func (s printArgs) String() string {
	return s.s
}
