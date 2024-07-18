package consumer

import (
	"github.com/IBM/sarama"
)

type consumer struct {
	client         client
	messageHandler *handler
	logger         sarama.StdLogger
}

func newConsumer(client client, messageHandler *handler, logger sarama.StdLogger) *consumer {
	return &consumer{
		client:         client,
		messageHandler: messageHandler,
		logger:         logger,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (r *consumer) Setup(sarama.ConsumerGroupSession) error {
	r.logger.Printf("consumer: setup...")
	// add call to dispatch a "setup" call to the client
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (r *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	r.logger.Printf("consumer: cleanup...")
	// add call to dispatch a "cleanup" call to the client
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (r *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case msg, msgChannelOk := <-claim.Messages():
			if !msgChannelOk {
				r.logger.Printf("consumer: message channel is closed")
				return errClosedMessageChannel
			}
			r.logger.Printf(
				"consumer: message received Ok: timestamp=%v, topic=%s, partition=%d, offset=%d",
				msg.Timestamp, msg.Topic, msg.Partition, msg.Offset,
			)

			// dispatch the message
			if err := r.messageHandler.Dispatch(session.Context(), msg); err != nil {
				r.logger.Printf("consumer: failed to dispatch message to client handler: err='%s'", err.Error())
				return newDispatchHandlerError(err)
			}

			// otherwise, we can commit this message now
			r.client.CommitMessage(session, msg)

		case <-session.Context().Done():
			// Should return when `session.Context()` is done.
			// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
			// https://github.com/IBM/sarama/issues/1192
			r.logger.Printf("consumer: context is Done. Exiting...")
			return nil
		}
	}
}
