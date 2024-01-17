package log

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type LoggerField struct {
	impl *zerolog.Event
}

func newLoggerField(impl *zerolog.Event) *LoggerField {
	return &LoggerField{impl: impl}
}

func SubDoc() *LoggerField {
	subDoc := zerolog.Dict()
	return newLoggerField(subDoc)
}

func (lf *LoggerField) Str(key string, val string) *LoggerField {
	lf.impl = lf.impl.Str(key, val)
	return lf
}

func (lf *LoggerField) Int(key string, val int) *LoggerField {
	lf.impl = lf.impl.Int(key, val)
	return lf
}

func (lf *LoggerField) Bool(key string, b bool) *LoggerField {
	lf.impl = lf.impl.Bool(key, b)
	return lf
}

func (lf *LoggerField) Duration(key string, d time.Duration) *LoggerField {
	lf.impl = lf.impl.Dur(key, d)
	return lf
}

func (lf *LoggerField) IPAddr(key string, ip net.IP) *LoggerField {
	lf.impl = lf.impl.IPAddr(key, ip)
	return lf
}

func (lf *LoggerField) UUID(key string, uuid uuid.UUID) *LoggerField {
	lf.impl = lf.impl.Str(key, uuid.String())
	return lf
}

func (lf *LoggerField) Properties(props *LoggerField) *LoggerField {
	lf.impl = lf.impl.Dict("properties", props.impl)
	return lf
}

// Deprecated: LegacyFields to support glamplify interface.
func (lf *LoggerField) LegacyFields(key string, f Fields) *LoggerField {
	if len(f) > 0 {
		lf.impl = lf.impl.Interface(key, f)
	}
	return lf
}

func (lf *LoggerField) Details(details string) {
	lf.impl.Msg(details)
}

func (lf *LoggerField) Detailsf(format string, v ...interface{}) {
	lf.impl.Msgf(format, v...)
}

func (lf *LoggerField) Send() {
	lf.impl.Send()
}

func (lf *LoggerField) doc(key string, props *LoggerField) *LoggerField {
	lf.impl = lf.impl.Dict(key, props.impl)
	return lf
}

func (lf *LoggerField) withFullStack() *LoggerField {
	lf.impl = lf.impl.Stack()
	return lf
}
