package consumer

import (
	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
)

type consumer struct {
	client         kafkaClient
	batchSize      int
	messageHandler handler
	logger         sarama.StdLogger
}

func newConsumer(client kafkaClient, batchSize int, messageHandler handler, logger sarama.StdLogger) *consumer {
	return &consumer{
		client:         client,
		batchSize:      batchSize,
		messageHandler: messageHandler,
		logger:         logger,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	c.logger.Printf("consumer: setup...")
	// add call to dispatch a "setup" call to the client
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	c.logger.Printf("consumer: cleanup...")
	// add call to dispatch a "cleanup" call to the client
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing loop and exit.
func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29

	var batch []*sarama.ConsumerMessage
	for {
		topic := claim.Topic()

		msg, err := c.getNextMessage(session, claim)
		if err != nil {
			// Either the session has closed or something went wrong reading the latest message.
			// dispatch what we have to the client and return the error
			if e := c.processBatch(session, topic, batch); e != nil {
				return errors.Errorf("consumer[%s]: failed to dispatch message to client handler: err='%s'", topic, e.Error())
			}
			return err
		}

		// successfully read a message, so add it to the batch
		batch = append(batch, msg)

		// if the batch is full, dispatch it
		if len(batch) >= c.batchSize {
			if err := c.processBatch(session, topic, batch); err != nil {
				// dispatch failed, so stop processing and return the error
				return err
			}
			batch = nil
		}
	}
}

func (c *consumer) getNextMessage(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (*sarama.ConsumerMessage, error) {
	topic := claim.Topic()
	select {
	case msg, msgChannelOk := <-claim.Messages():
		if !msgChannelOk {
			c.logger.Printf("consumer[%s]: message channel is closed", topic)
			return nil, errClosedMessageChannel
		}

		c.logger.Printf(
			"consumer[%s]: message received Ok: timestamp=%v, partition=%d, offset=%d",
			topic, msg.Timestamp, msg.Partition, msg.Offset,
		)
		return msg, nil

	case <-session.Context().Done():
		c.logger.Printf("consumer[%s]: context is Done. Exiting...", topic)
		return nil, errDoneMessageChannel
	}
}

func (c *consumer) processBatch(session sarama.ConsumerGroupSession, topic string, batch []*sarama.ConsumerMessage) error {
	err := c.dispatchBatch(session, topic, batch)
	c.logger.Printf("consumer[%s]: committing batch offset", topic)
	c.client.Commit(session)
	return err
}

func (c *consumer) dispatchBatch(session sarama.ConsumerGroupSession, topic string, batch []*sarama.ConsumerMessage) error {
	for _, msg := range batch {
		if err := c.messageHandler.Dispatch(session.Context(), msg); err != nil {
			c.logger.Printf("consumer[%s]: failed to dispatch message to client handler: err='%s'", topic, err.Error())
			return newDispatchHandlerError(topic, err)
		}

		c.logger.Printf("consumer[%s]: committing message with offset=%d", topic, msg.Offset)
		c.client.CommitMessage(session, msg)
	}
	return nil
}
