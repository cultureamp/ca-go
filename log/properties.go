package log

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Property contains an element of the log, usually a key-value pair.
type Property struct {
	impl *zerolog.Event
}

func newLoggerField(impl *zerolog.Event) *Property {
	return &Property{impl: impl}
}

// SubDoc creates a new sub-document with in list of log properties.
func SubDoc() *Property {
	subDoc := zerolog.Dict()
	return newLoggerField(subDoc)
}

// Str adds the property key with val as a string to the log.
func (lf *Property) Str(key string, val string) *Property {
	lf.impl = lf.impl.Str(key, val)
	return lf
}

// Int adds the property key with val as an int to the log.
func (lf *Property) Int(key string, val int) *Property {
	lf.impl = lf.impl.Int(key, val)
	return lf
}

// Bool adds the property key with val as an bool to the log.
func (lf *Property) Bool(key string, b bool) *Property {
	lf.impl = lf.impl.Bool(key, b)
	return lf
}

// Duration adds the property key with val as an time.Duration to the log.
func (lf *Property) Duration(key string, d time.Duration) *Property {
	lf.impl = lf.impl.Dur(key, d)
	return lf
}

// Time adds the property key with val as an uuid.UUID to the log.
func (lf *Property) Time(key string, t time.Time) *Property {
	// uses zerolog.TimeFieldFormat which we set to time.RFC3339
	lf.impl = lf.impl.Time(key, t)
	return lf
}

// IPAddr adds the property key with val as an net.IP to the log.
func (lf *Property) IPAddr(key string, ip net.IP) *Property {
	lf.impl = lf.impl.IPAddr(key, ip)
	return lf
}

// UUID adds the property key with val as an uuid.UUID to the log.
func (lf *Property) UUID(key string, uuid uuid.UUID) *Property {
	lf.impl = lf.impl.Str(key, uuid.String())
	return lf
}

// Properties adds an entire sub-document of type Property to the log.
func (lf *Property) Properties(props *Property) *Property {
	lf.impl = lf.impl.Dict("properties", props.impl)
	return lf
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

func (lf *Property) doc(key string, props *Property) *Property {
	lf.impl = lf.impl.Dict(key, props.impl)
	return lf
}

func (lf *Property) withFullStack() *Property {
	lf.impl = lf.impl.Stack()
	return lf
}
