package kms

import (
	"context"
	b64 "encoding/base64"

	"github.com/cultureamp/ca-go/x/awsconfig"

	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/pkg/errors"
)

type KMSClient interface {
	Encrypt(ctx context.Context, params *awskms.EncryptInput, optFns ...func(*awskms.Options)) (*awskms.EncryptOutput, error)
	Decrypt(ctx context.Context, params *awskms.DecryptInput, optFns ...func(*awskms.Options)) (*awskms.DecryptOutput, error)
}

type KMSEncrypt interface {
	Encrypt(ctx context.Context, plainStr string) (encryptedStr *string, err error)
	Decrypt(ctx context.Context, encryptedStr string) (decryptedStr *string, err error)
}

type kmsEncrypt struct {
	client KMSClient
	keyID  *string
}

func NewKMSWithClient(keyID string, client KMSClient) KMSEncrypt {
	return &kmsEncrypt{client, &keyID}
}

func NewKMS(keyID string) KMSEncrypt {
	return &kmsEncrypt{nil, &keyID}
}

func (k *kmsEncrypt) Encrypt(ctx context.Context, plainStr string) (*string, error) {
	if k.client == nil {
		svc, err := k.getNewServiceClient(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get kms client")
		}
		k.client = svc
	}

	input := &awskms.EncryptInput{
		KeyId:     k.keyID,
		Plaintext: []byte(plainStr),
	}

	result, err := k.client.Encrypt(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt with kms")
	}

	blobString := b64.StdEncoding.EncodeToString(result.CiphertextBlob)
	return &blobString, nil
}

func (k *kmsEncrypt) Decrypt(ctx context.Context, encryptedStr string) (*string, error) {
	if k.client == nil {
		svc, err := k.getNewServiceClient(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get kms client")
		}
		k.client = svc
	}

	blob, err := b64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	input := &awskms.DecryptInput{
		CiphertextBlob: blob,
	}

	result, err := k.client.Decrypt(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt with kms")
	}

	decStr := string(result.Plaintext)

	return &decStr, nil
}

func (k *kmsEncrypt) getNewServiceClient(ctx context.Context) (*awskms.Client, error) {
	cfg, err := awsconfig.GetAwsConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get aws config")
	}

	svc := awskms.NewFromConfig(cfg)

	return svc, nil
}
