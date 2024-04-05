# ca-go/launchdarkly
Package flags provides access to feature flags and product toggles. It wraps
the LaunchDarkly SDK and exposes a convenient and consistent way of configuring
and using the client.

The client is configured automatically based on the presence of the
LAUNCHDARKLY_CONFIGURATION environment variable which contains a JSON structured
string. You should declare this variable in your CDK configuration for your
infrastructure. The correct value for the environment your service is running
in can be retrieved from the AWS Secrets Manager under the key
`/common/launchdarkly-ops/sdk-configuration/<farm>`.

You can provide overrides for some of the properties of LAUNCHDARKLY_CONFIGURATION.
Refer to the documentation for the WithProxyMode and WithLambdaMode options in
config.go.

If the LAUNCHDARKLY_CONFIGURATION variable does not exist, the SDK will fall-back
to test mode. Test mode disables connections to LaunchDarkly and allows you to
specify your own values for flags. By default, it attempts to find a file named
.ld-flags.json in the directory that you invoked your Go program from. Failing
this, it configures the SDK to use dynamic test data sourced at runtime. You can
also specify your own path to a JSON file to source flag data from. See
WithTestMode() in config.go for more information.

The client can be configured and used as a managed singleton or as an
instance returned from a constructor function. The managed singleton provides
a layer of convenience by removing the need for your application to maintain
a handle on the flags client.

To configure the client as a singleton:

	err := flags.Configure()
	if err != nil {
	  // handle invalid configuration
	}

To configure the client as a instance that you manage:

	client, err := flags.NewClient()
	if err != nil {
	  // handle invalid configuration
	}

	err = client.Connect()
	if err != nil {
	  // handle errors connecting to LaunchDarkly
	}

The client will attempt to proxy requests through the LD Relay by default. You
can optionally choose to connect directly to DynamoDB by specifying the
WithLambdaMode() option to the flags.NewClient() or flags.Configure() functions.

Querying for flags is done on the client instance. You can get instance from the
managed singleton with GetDefaultClient():

	client, err := flags.GetDefaultClient()
	if err != nil {
	  // client not configured or connected
	}

A typical query takes three pieces of data:
 1. The flag name (the "key" within the LaunchDarkly UI).
 2. The evaluation context, which contains the identifiers and attributes of an
    entity that you wish to query the state of a flag for. See the
    evaluationcontext package for more information.
 3. The fallback value to return if an evaluation error occurs. This value will
    always be reflected as the value of the flag if err is not nil.

In most cases, the client can automatically build the evaluation context from
the request context (provided the context has been augmented with the
ca-go/request package):

	flagVal, err := client.QueryBool(ctx, "my-flag", false)

You can also supply your own evaluation context:

	evalcontext := flags.NewEvaluationContext(
	          flags.WithUserID("user-id"),
	          flags.WithUserAccountID("account-id"),
	)

	val, err := client.QueryBoolWithEvaluationContext("my-flag", evalcontext, false)

You will not need to manually shut down your SDK in most situations. If you
know your application is about to terminate, or if you're testing an app,
you should manually Shutdown() the LaunchDarkly client before quitting to ensure
it delivers any pending analytics events to LaunchDarkly:

	client.Shutdown()

## evaluation context

Package evaluationcontext defines the attributes for contexts you can
provide in a query for a flag or product toggle. Constructor functions are
exposed to create valid instances of evaluation contexts.

An evaluation context is simply a bag of attributes keyed by a unique
identifier. These values are used in two places:
 1. When creating a flag or segment, the attributes are used to form the
    targeting rules. For example, "if the user's realUserID is 123, return
    false for this flag".
 2. When querying for a flag in your service, the SDK uses the attributes
    to evaluate the rules for the flag to return the correct value. Using
    the same example as above, supplying an EvaluationContext with a realUserID
    of 123 would cause the flag to evaluate to false.

You should always supply as many attributes as you can to give yourself more
flexibility when writing new targeting rules. When you query a flag containing a
rule that works on attribute "foo", you must supply attribute "foo" in the
evaluation context.

User and Survey are now deprecated but have not been removed for backwards
compatibility. In order to upgrade use EvaluationContext instead there are
three steps to follow:
 1. Install the latest version of this package.
 2. Update targeting rules in LD to use contexts as explained on [confluence].
 3. Change all uses of User and Survey helper functions to EvaluationContext.

[confluence]: https://cultureamp.atlassian.net/wiki/spaces/CoSe/pages/2820768174/Guide+to+flag+targeting+in+LaunchDarkly
