package consumer

import (
	"encoding/binary"

	"github.com/go-errors/errors"
)

const (
	schemaBytes = 5 // AvroMagicByte is the magic byte used by Confluent Schema Registry.
)

type decoder interface {
	Decode(value []byte) (string, error)
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
func (d *avroDecoder) Decode(value []byte) (string, error) {
	value, err := d.decodeAsBytes(value)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func (d *avroDecoder) decodeAsBytes(value []byte) ([]byte, error) {
	if len(value) < schemaBytes {
		return nil, errors.Errorf("failed to decode: message missing schema id")
	}

	// Recover the schema id from the message and use the
	// client to retrieve the schema from Schema Registry.
	// Then use it to deserialize the record accordingly.
	schemaID := binary.BigEndian.Uint32(value[1:5])
	schema, err := d.client.GetSchemaByID(int(schemaID))
	if err != nil {
		return nil, err
	}

	codec := schema.Codec()
	native, _, err := codec.NativeFromBinary(value[5:])
	if err != nil {
		return nil, err
	}

	text, err := codec.TextualFromNative(nil, native)
	if err != nil {
		return nil, err
	}

	return text, nil
}
