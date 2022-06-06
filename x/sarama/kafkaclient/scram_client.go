package kafkaclient

import (
	"crypto/sha256"
	"crypto/sha512"

	"github.com/xdg-go/scram"
)

var (
	SHA256 scram.HashGeneratorFcn = sha256.New
	SHA512 scram.HashGeneratorFcn = sha512.New
)

// Code originally sourced from https://github.com/Shopify/sarama/blob/4c0bbf8d8ee3a91ee008efe5daa9d4abc42c10fa/examples/sasl_scram_client/scram_client.go
// See LICENSE file for licensing details.

// xDGSCRAMClient enables the use of the SHA256 and SHA512 SCRAM authentication
// mechanisms with the Sarama kafka client library. If no generator is specified,
// SHA512 is used.
type xDGSCRAMClient struct {
	client             *scram.Client
	clientConversation *scram.ClientConversation
	HashGenerator      scram.HashGeneratorFcn
}

func (x *xDGSCRAMClient) Begin(userName, password, authzID string) error {
	var err error
	if x.HashGenerator == nil {
		// SHA512 is our default setting
		x.HashGenerator = SHA512
	}

	x.client, err = x.HashGenerator.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}

	x.clientConversation = x.client.NewConversation()

	return nil
}

func (x *xDGSCRAMClient) Step(challenge string) (string, error) {
	return x.clientConversation.Step(challenge)
}

func (x *xDGSCRAMClient) Done() bool {
	return x.clientConversation.Done()
}
