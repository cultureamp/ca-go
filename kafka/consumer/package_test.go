package consumer

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPackage(t *testing.T) {
	ctx := context.Background()
	ch, stop := Consume(ctx, "topic-name")
	// when finished close it
	defer func() {
		err := stop()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	// Read the next message from the topic
	msg, ok := <-ch
	fmt.Printf("Channel open=%t, message=%v\n", ok, msg)
}

func TestPackageWithTimeout(t *testing.T) {
	ctx := context.Background()
	ch, stop := Consume(ctx, "topic-name")
	// when finished close it
	defer func() {
		err := stop()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	select {
	case msg, ok := <-ch:
		fmt.Printf("Channel open=%t, message=%v\n", ok, msg)

	case <-time.After(time.Duration(1) * time.Second):
		fmt.Println("No message received before timeout")
	}
}

func TestPackageWithDeadline(t *testing.T) {
	deadline := time.Duration(1) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	ch, stop := Consume(ctx, "employee.v1")
	// when finished close it
	defer func() {
		err := stop()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

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
}
