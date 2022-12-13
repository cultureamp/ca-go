// Package consumer provides a high level API for consuming messages from a Kafka
// topic.
//
// Consumer
//
// A single consumer is the most basic use case when consuming messages from a
// topic.
//
// c := consumer.NewConsumer(dialer, config)
// err := c.Run(ctx, handler)
//
// Consumer Group
//
// Multiple consumers can be created as Group to consume messages at a higher rate.
//
// g := consumer.Group(dialer, config)
// errCh := g.Run(ctx, handler)
// for err := range errCh {
//     if err != nil {
// 	       panic(err)
//     }
// }
//
// The consumer count and for the group is specified in the config parameter, which
// determines the number of goroutines to spawn. Each goroutine has their own
// Consumer and is set up so that individual consumer errors are reported back
// via the error channel returned from Group.Run.
//
// It is important to note that Kafka is ultimately responsible for managing group
// members. As a result, a consumer group can easily be spread across multiple
// instances by simply using the same group ID for each Group. Kafka will then
// take care of re-balancing the group if members are added/removed.
//
// Consumers in this package also do not support publishing failed messages to a
// dead letter queue out of the box. If this behaviour is required, it can simply
// be implemented in your Handler if an error occurs. However, it is important to
// return nil and not the error in this Handler approach so that the consumer can
// safely proceed to the next message.
package consumer
