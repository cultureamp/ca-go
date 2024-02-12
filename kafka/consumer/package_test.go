package consumer

import (
	"fmt"
	"testing"
)

func TestPackage(t *testing.T) {
	ch := Run("topic-name")

	// Read the next message from the topic
	msg, ok := <-ch
	fmt.Printf("Channel open=%t, message=%v\n", ok, msg)

	/*
		open := true
		// read next 10 messages
		for i := 1; i < 10 && open; i++ {
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

		// typical use would be read until channel is closed
		open = true
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
