package consumer

import (
	avro "github.com/riferrei/srclient"
)

type schemaRegistryClient interface {
	GetSchemaByID(id int) (*avro.Schema, error)
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

func (c *avroSchemaRegistryClient) GetSchemaByID(id int) (*avro.Schema, error) {
	return c.client.GetSchema(id)
}
