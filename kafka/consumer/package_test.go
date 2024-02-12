package consumer

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPackage(t *testing.T) {
	ctx := context.Background()
	ch := Run(ctx, "topic-name")

	// Read the next message from the topic
	msg, ok := <-ch
	fmt.Printf("Channel open=%t, message=%v\n", ok, msg)

	/*
		// typical use would be read until channel is closed
		open := true
		for open {
			select {
			case msg, ok := <-ch:
				if ok {
					fmt.Println(msg)
				} else {
					fmt.Println("topic closed")
					open = false
				}
			}
		}
	*/
}

func TestPackageWithTimeout(t *testing.T) {
	ctx := context.Background()
	ch := Run(ctx, "topic-name")

	select {
	case msg, ok := <-ch:
		fmt.Printf("Channel open=%t, message=%v\n", ok, msg)

	case <-time.After(time.Duration(30) * time.Second):
		fmt.Println("No message received before timeout")
	}
}

func TestPackageWithDeadline(t *testing.T) {
	deadline := time.Duration(30) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	ch := Run(ctx, "topic-name")

	select {
	case msg, ok := <-ch:
		fmt.Printf("Channel open=%t, message=%v\n", ok, msg)

	case <-ctx.Done():
		fmt.Println("No message received before context deadline")
	}
}
