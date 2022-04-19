// Package flags provides access to feature flags and product toggles. It wraps
// the LaunchDarkly SDK and exposes a convenient and consistent way of configuring
// and using the client.
//
// The client is configured automatically based on the presence of the
// LAUNCHDARKLY_CONFIGURATION environment variable which contains a JSON structured
// string. You should declare this variable in your CDK configuration for your
// infrastructure. The correct value for the environment your service is running
// in can be retrieved from the AWS Secrets Manager under the key
// `/common/launchdarkly-ops/sdk-configuration/<farm>`.
//
// The client can be configured and used as a managed singleton or as an
// instance returned from a constructor function. The managed singleton provides
// a layer of convenience by removing the need for your application to maintain
// a handle on the flags client.
//
// To configure the client as a singleton:
//   err := flags.Configure()
//   if err != nil {
//     // handle invalid configuration
//   }
//
//   err = flags.Connect()
//   if err != nil {
//     // handle errors connecting to LaunchDarkly
//   }
//
// To configure the client as a instance that you manage:
//   client, err := flags.NewClient()
//   if err != nil {
//     // handle invalid configuration
//   }
//
// The client will attempt to proxy requests through the LD Relay by default. You
// can optionally choose to connect directly to DynamoDB by specifying the
// WithLambdaMode() option to the flags.NewClient() or flags.Configure() functions.
//
// Querying for flags is done on the client instance. You can get instance from the
// managed singleton with GetDefaultClient():
//   client, err := flags.GetDefaultClient()
//   if err != nil {
//     // client not configured or connected
//   }
//
// A typical query takes three pieces of data:
// 1) The flag name (the "key" within the LaunchDarkly UI).
// 2) The evaluation context, which contains the identifiers and attributes of an
//    entity that you wish to query the state of a flag for. See the
//    evaluationcontext package for more information.
// 3) The fallback value to return if an evaluation error occurs. This value will
//    always be reflected as the value of the flag if err is not nil.
//
// In most cases, the client can automatically build the evaluation context from
// the request context (provided the context has been augmented with the
// ca-go/request package):
//   flagVal, err := client.QueryBool(ctx, "my-flag", false)
//
// You can also supply your own evaluation context:
//   user := flags.NewUser(
//             "user-id",
//             flags.WithAccountID("account-id"),
//   )
//
//   val, err := client.QueryBoolWithEvaluationContext("my-flag", user, false)
//
// When your application is shutting down, you should call Shutdown() to gracefully
// close connections to LaunchDarkly:
//   client.Shutdown()
package flags
