package consumer

import (
	"github.com/IBM/sarama"
)

type handler struct {
	client     kafkaClient
	decoder    decoder
	dispatcher dispatcher
	logger     sarama.StdLogger
}

func newHandler(client kafkaClient, decoder decoder, dispatcher dispatcher, logger sarama.StdLogger) *handler {
	return &handler{
		client:     client,
		dispatcher: dispatcher,
		decoder:    decoder,
		logger:     logger,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c *handler) Setup(sarama.ConsumerGroupSession) error {
	c.logger.Printf("consumer: setup...")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (c *handler) Cleanup(sarama.ConsumerGroupSession) error {
	c.logger.Printf("consumer: cleanup...")
	// We don't dispatch any outstanding messages in the batch to the client,
	// but neither will those be marked as committed
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing loop and exit.
func (c *handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29

	for {
		topic := claim.Topic()

		msg, err := c.getNextMessage(session, claim)
		if err != nil {
			// channel error (closed or context done?), so stop processing and return the error
			return err
		}

		if err := c.processMessage(session, topic, msg); err != nil {
			// dispatch failed, so stop processing and return the error
			return err
		}
	}
}

func (c *handler) getNextMessage(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (*sarama.ConsumerMessage, error) {
	topic := claim.Topic()
	select {
	case msg, msgChannelOk := <-claim.Messages():
		if !msgChannelOk {
			c.logger.Printf("consumer[%s]: message channel read error. Closed?", topic)
			return nil, errClosedMessageChannel
		}

		c.logger.Printf(
			"consumer[%s]: message received Ok: timestamp=%v, partition=%d, offset=%d",
			topic, msg.Timestamp, msg.Partition, msg.Offset,
		)
		return msg, nil

	case <-session.Context().Done():
		c.logger.Printf("acceptor[%s]: context is Done. Exiting...", topic)
		return nil, session.Context().Err()
	}
}

func (c *handler) processMessage(session sarama.ConsumerGroupSession, topic string, msg *sarama.ConsumerMessage) error {
	// 1. Decode the msg.Value []byte to a string using the avro decoder
	text, err := c.decoder.Decode(msg.Value)
	if err != nil {
		c.logger.Printf("consumer[%s]: failed to decode message[%s]: err='%s'", topic, string(msg.Key), err.Error())
		return newMessageHandlerError(topic, err)
	}

	message := &ReceivedMessage{
		Timestamp: msg.Timestamp,
		Topic:     msg.Topic,
		Offset:    msg.Offset,
		Key:       msg.Key,
		Value:     text,
	}

	// 2. Dispatch the message to the client's handler
	if err := c.dispatcher.Dispatch(session.Context(), message); err != nil {
		c.logger.Printf("consumer[%s]: failed to dispatch message[%s] to client handler: err='%s'", topic, string(msg.Key), err.Error())
		return newMessageHandlerError(topic, err)
	}

	// 3. Mark the message as successfully consumed
	c.logger.Printf("consumer[%s]: marking message[%s] as successfully consumed", topic, string(msg.Key))
	c.client.MarkMessageConsumed(session, msg, "done")

	// Given https://medium.com/@moabbas.ch/effective-kafka-consumption-in-golang-a-comprehensive-guide-aac54b5b79f0
	// It looks like calling MarkMessageConsumed() is enough to commit the message offset.
	// If not, then uncomment the lines below.
	//
	// c.logger.Printf("acceptor[%s]: committing message offset[%d]", topic, msg.Offset)
	// c.client.Commit(session)
	return nil
}
