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

	props := Add().
		Str("str", "value").
		Int("int", 1).
		Bool("bool", true).
		Duration("dur", duration).
		Time("time", then).
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
