package cago

import (
	"context"

	"github.com/cultureamp/ca-go/log"
)



func test() {
	ctx := context.Background()

	log.Debug(ctx, "hello")
}