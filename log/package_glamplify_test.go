package log_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/log"
)

func TestCommonGlamplifyExamples(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	f := log.Fields{
		"key1":    "string value",
		"key2":    1,
		"now":     now.Format(time.RFC3339),
		"later":   time.Now(),
		"details": "detailed message",
	}
	log.LogDebug(ctx, "log_fields", f)
	log.LogInfo(ctx, "log_fields", f)
	log.LogWarn(ctx, "log_fields", f)
	log.LogError(ctx, "log_fields", errors.New("test error"), f)

	// log.LogFatal calls os.exit() so this is hard to test!

	defer recoverFromPanic()
	log.LogPanic(ctx, "panic_error", errors.New("test error"), f)
}

func recoverFromLogPanic() {
	if saved := recover(); saved != nil {
		// convert to an error if it's not one already
		err, ok := saved.(error)
		if !ok {
			err = errors.New(fmt.Sprint(saved))
		}

		log.Error("recovered_from_panic", err).Send()
	}
}
