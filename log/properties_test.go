package log

import (
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func TestFieldTypes(t *testing.T) {
	var ipv4 net.IP

	then := time.Now()
	u := uuid.New()
	duration := time.Since(then)
	f := func(e *zerolog.Event) { e.Str("func", "val") }
	var i64 int64
	var ui64 uint64

	i64 = 123
	ui64 = 123

	props := Add().
		Str("str", "value").
		Int("int", 1).
		Int64("int64", i64).
		UInt64("uint64", ui64).
		Bool("bool", true).
		Duration("dur", duration).
		Time("time", then).
		IPAddr("ipaddr", ipv4).
		UUID("uuid", u).
		Func(f)

	Debug("debug_with_all_field_types").
		WithRequestTracing(nil).
		Properties(props).
		Details("logging should contain all types")

	Debug("debug_with_all_field_types").
		WithRequestTracing(nil).
		Properties(props).
		Detailsf("logging should contain all types: %s", "ok")

	Debug("debug_with_all_field_types").
		WithRequestTracing(nil).
		Properties(props).
		Send()
}
