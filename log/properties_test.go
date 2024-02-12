package log

import (
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestFieldTypes(t *testing.T) {
	var ipv4 net.IP

	then := time.Now()
	u := uuid.New()
	duration := time.Since(then)

	var i64 int64
	var u64 uint64

	i64 = -212123121231323434
	u64 = 345323423423423434
	props := SubDoc().
		Str("str", "value").
		Int("int", 1).
		Int64("int64", i64).
		UInt64("uint64", u64).
		Bool("bool", true).
		Duration("dur", duration).
		IPAddr("ipaddr", ipv4).
		UUID("uuid", u)

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
