package consumer

import (
	"encoding/binary"

	"github.com/IBM/sarama"
	avro "github.com/riferrei/srclient"
)

type decoder interface {
	Decode(msg *sarama.ConsumerMessage) ([]byte, error)
	DecodeAsString(msg *sarama.ConsumerMessage) (string, error)
}

type avroSchemaRegistryClient struct {
	client *avro.SchemaRegistryClient
}

func newAvroSchemaRegistryClient(schemaRegistryURL string) *avroSchemaRegistryClient {
	client := avro.CreateSchemaRegistryClient(schemaRegistryURL)

	return &avroSchemaRegistryClient{
		client: client,
	}
}

func (c *avroSchemaRegistryClient) Decode(msg *sarama.ConsumerMessage) ([]byte, error) {
	// Recover the schema id from the message and use the
	// client to retrieve the schema from Schema Registry.
	// Then use it to deserialize the record accordingly.
	schemaID := binary.BigEndian.Uint32(msg.Value[1:5])
	schema, err := c.client.GetSchema(int(schemaID))
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

func (c *avroSchemaRegistryClient) DecodeAsString(msg *sarama.ConsumerMessage) (string, error) {
	value, err := c.Decode(msg)
	if err != nil {
		return "", err
	}

	return string(value), nil
}
