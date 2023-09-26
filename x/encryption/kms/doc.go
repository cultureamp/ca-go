// package kms provides the ability to encrypt/decrypt data based on AWS KMS
//
// # To create the kms client
//
// ```
// cfg, err := config.LoadDefaultConfig(
//
//		ctx,
//		config.WithRegion(settings.AwsRegion),
//	)
//
// client := awskms.NewFromConfig(cfg)
// ```
//
// then create an encryptor
//
// `encryptor := kms.NewEncryptor(settings.DataKMSKeyID, client)`
//
// then you can encrypt your data by
//
// `encryptor.Encrypt(ctx, metadata.MergeAccountToken)`
//
// then you can decrypt your data by
//
// `encryptor.Decrypt(ctx, metadataItem.MergeAccountToken)`
package kms
