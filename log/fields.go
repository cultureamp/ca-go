package log

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Field contains an element of the log, usually a key-value pair.
type Field struct {
	impl *zerolog.Event
}

func newLoggerField(impl *zerolog.Event) *Field {
	return &Field{impl: impl}
}

// Add creates a new custom log properties list.
func Add() *Field {
	subDoc := zerolog.Dict()
	return newLoggerField(subDoc)
}

// Str adds the property key with val as a string to the log.
// Note: Empty string values will not be logged.
func (lf *Field) Str(key string, val string) *Field {
	if val == "" {
		return lf
	}

	lf.impl = lf.impl.Str(key, val)
	return lf
}

// Int adds the property key with val as an int to the log.
func (lf *Field) Int(key string, val int) *Field {
	lf.impl = lf.impl.Int(key, val)
	return lf
}

// Int64 adds the property key with val as an int64 to the log.
func (lf *Field) Int64(key string, val int64) *Field {
	lf.impl = lf.impl.Int64(key, val)
	return lf
}

// UInt64 adds the property key with val as an uint64 to the log.
func (lf *Field) UInt64(key string, val uint64) *Field {
	lf.impl = lf.impl.Uint64(key, val)
	return lf
}

// Bool adds the property key with val as an bool to the log.
func (lf *Field) Bool(key string, b bool) *Field {
	lf.impl = lf.impl.Bool(key, b)
	return lf
}

// Duration adds the property key with val as an time.Duration to the log.
func (lf *Field) Duration(key string, d time.Duration) *Field {
	lf.impl = lf.impl.Dur(key, d)
	return lf
}

// Time adds the property key with val as an uuid.UUID to the log.
func (lf *Field) Time(key string, t time.Time) *Field {
	// uses zerolog.TimeFieldFormat which we set to time.RFC3339
	lf.impl = lf.impl.Time(key, t)
	return lf
}

// IPAddr adds the property key with val as an net.IP to the log.
func (lf *Field) IPAddr(key string, ip net.IP) *Field {
	lf.impl = lf.impl.IPAddr(key, ip)
	return lf
}

// UUID adds the property key with val as an uuid.UUID to the log.
func (lf *Field) UUID(key string, uuid uuid.UUID) *Field {
	lf.impl = lf.impl.Str(key, uuid.String())
	return lf
}
