package consumer

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPackage(t *testing.T) {
	ctx := context.Background()
	ch := Consume(ctx, "topic-name")

	// Read the next message from the topic
	msg, ok := <-ch
	fmt.Printf("Channel open=%t, message=%v\n", ok, msg)

	// when finished close it
	err := Stop("topic-name")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func TestPackageWithTimeout(t *testing.T) {
	ctx := context.Background()
	ch := Consume(ctx, "topic-name")

	select {
	case msg, ok := <-ch:
		fmt.Printf("Channel open=%t, message=%v\n", ok, msg)

	case <-time.After(time.Duration(1) * time.Second):
		fmt.Println("No message received before timeout")
	}

	// when finished close it
	err := Stop("topic-name")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func TestPackageWithDeadline(t *testing.T) {
	deadline := time.Duration(1) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	ch := Consume(ctx, "topic-name")

	ok := true
	for ok {
		select {
		case msg, ok := <-ch:
			fmt.Printf("Channel open=%t, message=%v\n", ok, msg)

		case <-ctx.Done():
			fmt.Println("Context deadline received. Stopping.")
			ok = false
		}
	}

	// when finished close it
	err := Stop("topic-name")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
