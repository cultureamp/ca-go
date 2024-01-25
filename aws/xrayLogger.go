package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/cultureamp/glamplify/log"
)

type xrayLogger struct {
	log *log.Logger
}

func newXrayLogger(ctx context.Context) *xrayLogger {
	logger := log.NewFromCtx(ctx)
	return &xrayLogger{
		log: logger,
	}
}

func (logger xrayLogger) Log(level xraylog.LogLevel, msg fmt.Stringer) {
	segment := logger.log.Event("xray_diagnostic")

	switch level {
	case xraylog.LogLevelDebug: segment.Debug(msg.String())
	case xraylog.LogLevelInfo: segment.Info(msg.String())
	case xraylog.LogLevelWarn: segment.Warn(msg.String())
	case xraylog.LogLevelError: segment.Error(errors.New(msg.String()))
	default: segment.Debug(msg.String())
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
