package consumer

import (
	"log"
	"os"
)

type testLogger struct {
	l *log.Logger
}

func newTestLogger() *testLogger {
	l := log.New(os.Stdout, "[Kafka-Consumer] ", log.LstdFlags)
	return &testLogger{
		l: l,
	}
}

func (d *testLogger) Print(v ...interface{}) {
	log.Print(v...)
}

func (d *testLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (d *testLogger) Println(v ...interface{}) {
	log.Println(v...)
}
