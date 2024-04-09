package log

import (
	"github.com/rs/zerolog"
)

// timestampHook allows us to dynamically set the time as time.Now for production, or a static value for tests.
type timestampHook struct {
	config *Config
}

// Run implemented the Hook interface.
func (t *timestampHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	// Uses zerolog.TimeFieldFormat which we set to time.RFC3339
	e.Time(zerolog.TimestampFieldName, t.config.TimeNow())
}
