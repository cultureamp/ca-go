// package secrets adds helper methods to fetch secrets from the AWS secret
// manager for use in AWS lambdas.
//
// **Secrets should be fetched during the cold start of the lambda**
//
// To create the secret manager client:
//
// `secretClient, err := secrets.NewAWSSecretsClient(envConfig.AwsRegion)`
//
// then fetch the secret using Get:
//
// `secretValue, err = secretClient.Get("/your/secret-name")`
package secrets
