package log

import (
	"github.com/rs/zerolog"
)

// EventHook adds a field "event" to the log which is the past-tense verb of what just happened, in snake_case.
// Ideally, "event" should be in non-trivial services, the event should be namespaced, delimited by periods (.),
// where the last section is the verb. (e.g. jobs.survey_launch.succeeded).
// It just parses "message" to snake case for now.
type EventHook struct{}

func (h EventHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Str("event", ToSnakeCase(msg))
}
