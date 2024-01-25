package aws

import (
	"testing"

	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/cultureamp/ca-go/log"
	"github.com/stretchr/testify/assert"
)

func Test_New_TraceLogger(t *testing.T) {
	log := newXrayLogger()
	assert.NotNil(t, log)
}

func Test_TraceLogger_RealWorld_Log(t *testing.T) {
	t.Setenv(log.LogQuietModeEnv, "false")
	t.Setenv(log.LogLevelEnv, "DEBUG")

	log := newXrayLogger()
	log.Log(xraylog.LogLevelDebug, newPrintArgs("debug"))
	log.Log(xraylog.LogLevelInfo, newPrintArgs("info"))
	log.Log(xraylog.LogLevelWarn, newPrintArgs("warn"))
	log.Log(xraylog.LogLevelError, newPrintArgs("error"))
}
