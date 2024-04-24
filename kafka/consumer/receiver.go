package consumer

import (
	"github.com/IBM/sarama"
)

type receiver struct {
	handler messageHandler
}

func newReceiver(handler Handler) *receiver {
	return &receiver{
		handler: *newMessageHandler(handler),
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (r *receiver) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (r *receiver) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (r *receiver) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				sarama.Logger.Printf("message channel was closed\n")
				return nil
			}
			sarama.Logger.Printf("message received: timestamp=%v, topic=%s, partition=%d, offset=%d\n", msg.Timestamp, msg.Topic, msg.Partition, msg.Offset)

			// dispatch the message
			if err := r.handler.dispatch(session.Context(), msg); err != nil {
				sarama.Logger.Printf("failed to dispatch message to handler: %w\n", err)
				return err
			}

			// otherwise, we can commit this message now
			session.MarkMessage(msg, "")

		case <-session.Context().Done():
			// Should return when `session.Context()` is done.
			// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
			// https://github.com/IBM/sarama/issues/1192
			return nil
		}
	}
}
