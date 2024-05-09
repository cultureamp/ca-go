package sentry

import "github.com/getsentry/sentry-go"

// RootCauseAsTitle uses error message as custom exception type. The general errors like *erros.errorString can end up
// grouping errors and makes it harder for us to find the latest error in Sentry issues dashboard.
func RootCauseAsTitle(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
	for i, exception := range event.Exception {
		event.Exception[i].Type = exception.Value
	}
	return event
}
