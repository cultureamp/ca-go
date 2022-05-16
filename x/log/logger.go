package log

import (
	"context"
	"encoding/json"
	"github.com/cultureamp/ca-go/x/request"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	*logrus.Entry
}

var log = &logrus.Logger{
	Out:          os.Stdout,
	Formatter:    &logrus.JSONFormatter{PrettyPrint: true},
	Hooks:        make(logrus.LevelHooks),
	Level:        logrus.TraceLevel,
	ReportCaller: true,
}

// setupFormatter decides the output formatter based on the environment where the app is running on.
// It uses text formatter with color if you run the app locally,
// while using json formatter if it's running on the cloud.
func setupDefaultFormatter(config EnvConfig) {
	if config.Farm == "local" {
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
	} else {
		log.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: true,
		})
	}
}

func convertFields(fields any) map[string]any {
	var fieldsMap map[string]interface{}
	data, _ := json.Marshal(fields)
	err := json.Unmarshal(data, &fieldsMap)
	if err != nil {
		log.WithError(err).Error("failed to parse logger fields from context")
	}
	return fieldsMap
}

func newLogger(ctx context.Context) *Logger {
	entry := log.WithContext(ctx).WithTime(time.Now())

	config, ok := EnvConfigFromContext(ctx)
	if ok {
		setupDefaultFormatter(config)
		entry = entry.WithFields(convertFields(config))
	}

	reqIds, ok := request.RequestIDsFromContext(ctx)
	if ok {
		entry = entry.WithFields(convertFields(reqIds))
	}

	userIds, ok := request.AuthenticatedUserFromContext(ctx)
	if ok {
		entry = entry.WithFields(convertFields(userIds))
	}

	return &Logger{
		entry,
	}
}

// NewFromCtx creates a new logger from a context, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromCtx(ctx context.Context) *Logger {
	if !ContextHasEnvConfig(ctx) {
		// add default env config if not exists in the context
		ctx = ContextWithDefaultEnvConfig(ctx)
	}
	return newLogger(ctx)
}

// NewFromRequest creates a new logger from a http.Request, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromRequest(r *http.Request) *Logger {
	return NewFromCtx(r.Context())
}

// this is copied from Glamplify, should we keep the getErrorValues() method or
// just using the location info provided by logrus.WithError()?
// WithError provides location info like this:
//   "file": "/Users/aimee.geng/develop/fusionauth-hooks/src/extuserupdate/lib/main.go:28",
//  "func": "main.main",
func (logger *Logger) WithErrorTrace(err error) *Logger {
	df := newErrorValues()
	return &Logger{
		logger.WithFields(df.getErrorValues(err)).WithError(err),
	}
}

func (logger *Logger) WithDatadogHook() {

}
