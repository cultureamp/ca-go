// package secrets adds helper methods to fetch secrets from the AWS secret
// manager for use in AWS lambdas.
//
// **Secrets should be fetched during the cold start of the lambda**
//
// # To retrieve a secret use
//
// `secretValue, err = secrets.Get("/your/secret-name")`
package secrets
