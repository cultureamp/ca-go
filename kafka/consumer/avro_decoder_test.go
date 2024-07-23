package consumer

import (
	"encoding/binary"
	"encoding/json"
	"testing"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAvroDecodeErrorWhenMessageIsNil(t *testing.T) {
	mockSchemaRegistryClient := newMockSchemaRegistryClient()
	decoder := newAvroDecoder(mockSchemaRegistryClient)

	_, err := decoder.Decode(nil)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to decode: message is nil")

	mockSchemaRegistryClient.AssertExpectations(t)
}

func TestAvroDecodeErrorWhenSchemaNotFound(t *testing.T) {
	mockSchemaRegistryClient := newMockSchemaRegistryClient()
	decoder := newAvroDecoder(mockSchemaRegistryClient)

	mockSchemaRegistryClient.On("GetSchemaByID", mock.Anything).Return(nil, errors.Errorf("failed to get schema"))

	bs := make([]byte, 10)
	binary.BigEndian.PutUint32(bs[1:], 1234)

	msg := &sarama.ConsumerMessage{}
	msg.Value = bs

	_, err := decoder.Decode(msg)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to get schema")

	mockSchemaRegistryClient.AssertExpectations(t)
}

func TestAvroDecode(t *testing.T) {
	mockSchemaRegistryClient := newMockSchemaRegistryClient()
	decoder := newAvroDecoder(mockSchemaRegistryClient)

	schema := testDecoderSchema(t)
	mockSchemaRegistryClient.On("GetSchemaByID", mock.Anything).Return(schema, nil)

	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	td := testData{ID: 123, Name: "Gopher"}
	value, _ := json.Marshal(td)
	codec := schema.Codec()
	native, _, err := codec.NativeFromTextual(value)
	assert.Nil(t, err)
	valueBytes, err := codec.BinaryFromNative(nil, native)
	assert.Nil(t, err)

	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)

	msg := &sarama.ConsumerMessage{}
	msg.Value = recordValue

	decodedJSON, err := decoder.Decode(msg)
	assert.Nil(t, err)

	var result testData
	err = json.Unmarshal([]byte(decodedJSON), &result)
	assert.Nil(t, err)

	assert.Equal(t, result.ID, 123)
	assert.Equal(t, result.Name, "Gopher")

	mockSchemaRegistryClient.AssertExpectations(t)
}

func testDecoderSchema(t *testing.T) *srclient.Schema {
	codec, err := goavro.NewCodec(`
	{
  "type": "record",
  "name": "TestObject",
  "namespace": "ca.dataedu",
  "fields": [
    {
      "name": "id",
      "type": "int"
    },
    {
      "name": "name",
      "type": "string"
    }
  ]
}`)
	assert.Nil(t, err)

	schema, err := srclient.NewSchema(1, "TestObject", srclient.Avro, 1, nil, codec, nil)
	assert.Nil(t, err)

	return schema
}

type testData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
