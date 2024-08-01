package consumer

import (
	"crypto/sha256"
	"crypto/sha512"

	"github.com/xdg-go/scram"
)

var (
	sha256Fn scram.HashGeneratorFcn = sha256.New
	sha512Fn scram.HashGeneratorFcn = sha512.New
)

type scramClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func newScramClient(hashGeneratorFcn scram.HashGeneratorFcn) *scramClient {
	return &scramClient{
		HashGeneratorFcn: hashGeneratorFcn,
	}
}

func (x *scramClient) Begin(userName, password, authzID string) error {
	client, err := x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}

	x.Client = client
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *scramClient) Step(challenge string) (string, error) {
	return x.ClientConversation.Step(challenge)
}

func (x *scramClient) Done() bool {
	return x.ClientConversation.Done()
}
