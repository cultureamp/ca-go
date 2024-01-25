package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/cultureamp/ca-go/log"
)

type xrayLogger struct {
	log *log.Logger
}

func newXrayLogger() *xrayLogger {
	config := log.NewLoggerConfig()
	logger := log.NewLogger(config)
	return &xrayLogger{
		log: logger,
	}
}

func (xl xrayLogger) Log(level xraylog.LogLevel, msg fmt.Stringer) {
	props := log.SubDoc().Str("message", msg.String())

	switch level {
	case xraylog.LogLevelDebug:
		log.Debug("xray_diagnostic").Properties(props).Send()
	case xraylog.LogLevelInfo:
		log.Info("xray_diagnostic").Properties(props).Send()
	case xraylog.LogLevelWarn:
		log.Warn("xray_diagnostic").Properties(props).Send()
	case xraylog.LogLevelError:
		log.Error("xray_diagnostic", errors.New("xray_diagnostic")).Properties(props).Send()
	default:
		log.Debug("xray_diagnostic").Properties(props).Send()
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
