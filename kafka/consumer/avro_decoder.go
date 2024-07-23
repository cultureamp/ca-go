package consumer

import (
	"encoding/binary"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
)

type decoder interface {
	Decode(msg *sarama.ConsumerMessage) (string, error)
}

type avroDecoder struct {
	client schemaRegistryClient
}

func newAvroDecoder(client schemaRegistryClient) *avroDecoder {
	return &avroDecoder{
		client: client,
	}
}

// Decode takes a Kafka consumer message and avro decodes the value field as a string (usually json).
func (d *avroDecoder) Decode(msg *sarama.ConsumerMessage) (string, error) {
	value, err := d.decodeAsBytes(msg)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func (d *avroDecoder) decodeAsBytes(msg *sarama.ConsumerMessage) ([]byte, error) {
	if msg == nil {
		return nil, errors.Errorf("failed to decode: message is nil")
	}

	// Recover the schema id from the message and use the
	// client to retrieve the schema from Schema Registry.
	// Then use it to deserialize the record accordingly.
	schemaID := binary.BigEndian.Uint32(msg.Value[1:5])
	schema, err := d.client.GetSchemaByID(int(schemaID))
	if err != nil {
		return nil, err
	}

	codec := schema.Codec()
	native, _, err := codec.NativeFromBinary(msg.Value[5:])
	if err != nil {
		return nil, err
	}

	value, err := codec.TextualFromNative(nil, native)
	if err != nil {
		return nil, err
	}

	return value, nil
}
