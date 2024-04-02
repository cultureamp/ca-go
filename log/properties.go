package log

import (
	"github.com/rs/zerolog"
)

// Property contains an element of the log, usually a key-value pair.
type Property struct {
	impl *zerolog.Event
}

func newLoggerProperty(impl *zerolog.Event) *Property {
	return &Property{impl: impl}
}

// Properties adds an entire sub-document of type Property to the log.
func (lf *Property) Properties(fields *Field) *Property {
	return lf.doc("properties", fields)
}

// Details adds the property 'details' with the val as a string to the log.
// This is a terminating Property that signals that the log statement is complete
// and can now be sent to the output.
//
// NOTICE: once this method is called, the *Property should be disposed.
// Calling Details twice can have unexpected result.
func (lf *Property) Details(details string) {
	lf.impl.Msg(details)
}

// Detailsf adds the property 'details' with the format and args to the log.
// This is a terminating Property that signals that the log statement is complete
// and can now be sent to the output.
//
// NOTICE: once this method is called, the *Property should be disposed.
// Calling Detailsf twice can have unexpected result.
func (lf *Property) Detailsf(format string, v ...interface{}) {
	lf.impl.Msgf(format, v...)
}

// Send terminates the log and signals that it is now complete and can be
// sent to the output.
//
// NOTICE: once this method is called, the *Property should be disposed.
func (lf *Property) Send() {
	lf.impl.Send()
}

func (lf *Property) doc(key string, fields *Field) *Property {
	lf.impl = lf.impl.Dict(key, fields.impl)
	return lf
}
