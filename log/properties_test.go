package log

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func ExampleLogger_Debug() {
	then := time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC)
	u := uuid.MustParse("e5fa7acf-1846-41b4-a2ee-80ecd86fb060")
	duration := time.Second * 42
	f := func(e *zerolog.Event) { e.Str("func", "val") }
	b := []byte("some bytes")

	var ui uint
	var i64 int64
	var ui64 uint64
	var f32 float32
	var f64 float64

	ui = 234
	i64 = 123
	ui64 = 123
	f32 = 32.32
	f64 = 64.64

	config := getExampleLoggerConfig("DEBUG")
	logger := NewLogger(config)

	props := Add().
		Str("str", "value").
		Int("int", 1).
		UInt("uint", ui).
		Int64("int64", i64).
		UInt64("uint64", ui64).
		Float32("float32", f32).
		Float64("float64", f64).
		Bool("bool", true).
		Bytes("bytes", b).
		Duration("dur", duration).
		Time("time", then).
		IPAddr("ipaddr", net.IPv4bcast).
		UUID("uuid", u).
		Func(f)

	logger.Debug("debug_with_all_field_types").
		Properties(props).
		Detailsf("logging should contain all types: %s", "ok")

	// Output:
	// 2020-11-14T11:30:32Z DBG event="logging should contain all types: ok" app=logger-test app_version=1.0.0 aws_account_id=development aws_region=def event=debug_with_all_field_types farm=local product=cago properties={"bool":true,"bytes":"some bytes","dur":"PT42S","float32":32.32,"float64":64.64,"func":"val","int":1,"int64":123,"ipaddr":"255.255.255.255","str":"value","time":"2023-11-14T11:30:32Z","uint":234,"uint64":123,"uuid":"e5fa7acf-1846-41b4-a2ee-80ecd86fb060"}
}

func getExampleLoggerConfig(sev string) *Config {
	config := NewLoggerConfig()
	config.AppName = "logger-test"
	config.AwsRegion = "def"
	config.Product = "cago"
	config.LogLevel = sev
	config.Quiet = false
	config.ConsoleWriter = true
	config.ConsoleColour = false
	config.TimeNow = func() time.Time {
		return time.Date(2020, 11, 14, 11, 30, 32, 0, time.UTC)
	}
	return config
}
