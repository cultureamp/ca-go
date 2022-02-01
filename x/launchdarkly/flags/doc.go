// Package flags provides access to feature flags and product toggles. It wraps
// the LaunchDarkly SDK and exposes a convenient and consistent way of configuring
// and using the client.
//
// The client can be configured and used as a managed singleton or as an
// instance returned from a constructor function. The managed singleton provides
// a layer of convenience by removing the need for your application to maintain
// a handle on the flags client.
//
// To configure the client as a singleton:
//   err := flags.Configure(
//	   flags.WithSDKKey("foobar"),
//   )
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
//   client, err := flags.NewClient(
//	   flags.WithSDKKey("foobar"),
//   )
//   if err != nil {
//     // handle invalid configuration
//   }
//
// The client can be configured to use the Relay Proxy in either Daemon or Proxy
// modes. See WithDaemonMode and WithProxyMode for more information.
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
//   flagVal, err := client.QueryBool(ctx, "flag.my-flag", false)
//
//   toggleVal, err := client.QueryToggle(ctx, "toggle.my-toggle", false)
//
// You can also supply your own evaluation context:
//   user := flags.NewUser(
//             "user-id",
//             flags.WithAccountID("account-id"),
//   )
//
//   val, err := client.QueryBoolWithEvaluationContext("flag.my-flag", user, false)
//
// When your application is shutting down, you should call Shutdown() to gracefully
// close connections to LaunchDarkly:
//   client.Shutdown()
package flags
