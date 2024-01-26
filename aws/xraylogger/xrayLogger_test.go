package aws

import (
	"testing"

	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/cultureamp/ca-go/log"
	"github.com/stretchr/testify/assert"
)

func Test_New_TraceLogger(t *testing.T) {
	xl := newXrayLogger()
	assert.NotNil(t, xl)
}

func Test_TraceLogger_RealWorld_Log(t *testing.T) {
	t.Setenv(log.LogQuietModeEnv, "false")
	t.Setenv(log.LogLevelEnv, "DEBUG")

	xl := newXrayLogger()
	xl.Log(xraylog.LogLevelDebug, newPrintArgs("debug"))
	xl.Log(xraylog.LogLevelInfo, newPrintArgs("info"))
	xl.Log(xraylog.LogLevelWarn, newPrintArgs("warn"))
	xl.Log(xraylog.LogLevelError, newPrintArgs("error"))
}
