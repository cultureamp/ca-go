package consumer

import (
	"github.com/IBM/sarama"
)

type receiver struct {
	client  kafkaClient
	handler *messageHandler
	logger  sarama.StdLogger
}

func newReceiver(client kafkaClient, handler Handler, logger sarama.StdLogger) *receiver {
	return &receiver{
		client:  client,
		handler: newMessageHandler(handler),
		logger:  logger,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (r *receiver) Setup(sarama.ConsumerGroupSession) error {
	r.logf("receiver: setup...")
	// add call to dispatch a "setup" call to the client
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (r *receiver) Cleanup(sarama.ConsumerGroupSession) error {
	r.logf("receiver: cleanup...")
	// add call to dispatch a "cleanup" call to the client
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
		case msg, msgChannelOk := <-claim.Messages():
			if !msgChannelOk {
				r.logf("receiver: message channel was closed")
				return nil
			}
			r.logf(
				"receiver: message received Ok: timestamp=%v, topic=%s, partition=%d, offset=%d",
				msg.Timestamp, msg.Topic, msg.Partition, msg.Offset,
			)

			// dispatch the message
			if err := r.handler.dispatch(session.Context(), msg); err != nil {
				r.logf("receiver: failed to dispatch message to client handler: '%s'", err.Error())
				return err
			}

			// otherwise, we can commit this message now
			r.client.CommitMessage(session, msg)

		case <-session.Context().Done():
			// Should return when `session.Context()` is done.
			// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
			// https://github.com/IBM/sarama/issues/1192
			r.logf("receiver: context is Done. Exiting...")
			return nil
		}
	}
}

func (r *receiver) logf(format string, v ...interface{}) {
	if r.logger == nil {
		return
	}

	r.logger.Printf(format, v...)
}
