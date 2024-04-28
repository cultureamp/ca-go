package log

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
)

// Property contains an element of the log, usually a key-value pair.
type Property struct {
	impl      *zerolog.Event
	configErr error
}

func newLoggerProperty(impl *zerolog.Event, config *Config) *Property {
	// Default is to assume there is no config (ie. from mocks, tests, etc)
	var err error = errors.Errorf("missing logger config")
	if config != nil {
		err = config.isValid()
	}

	return &Property{
		impl:      impl,
		configErr: err,
	}
}

// Properties adds an entire sub-document of type Property to the log.
func (lf *Property) Properties(fields *Field) *Property {
	return lf.doc("properties", fields)
}

// Details adds the property 'details' with the val as a string to the log.
// This is a terminating Property that signals that the log statement is complete
// and can now be sent to the output. It returns nil on success, or an error if
// there was a problem.
//
// NOTICE: once this method is called, the *Property should be disposed.
// Calling Details twice can have unexpected results.
func (lf *Property) Details(details string) error {
	lf.impl.Msg(details)
	return lf.configErr
}

// Detailsf adds the property 'details' with the format and args to the log.
// This is a terminating Property that signals that the log statement is complete
// and can now be sent to the output.It returns nil on success, or an error if
// there was a problem.
//
// NOTICE: once this method is called, the *Property should be disposed.
// Calling Detailsf twice can have unexpected results.
func (lf *Property) Detailsf(format string, v ...interface{}) error {
	lf.impl.Msgf(format, v...)
	return lf.configErr
}

// Send terminates the log and signals that it is now complete and can be
// sent to the output. It returns nil on success, or an error if
// there was a problem.
//
// NOTICE: once this method is called, the *Property should be disposed.
// Calling Send twice can have unexpected results.
func (lf *Property) Send() error {
	lf.impl.Send()
	return lf.configErr
}

func (lf *Property) doc(key string, fields *Field) *Property {
	lf.impl = lf.impl.Dict(key, fields.impl)
	return lf
}
