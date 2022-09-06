module github.com/cultureamp/ca-go

go 1.18

require (
	github.com/aws/aws-sdk-go v1.44.91
	github.com/launchdarkly/go-server-sdk-dynamodb v1.1.1
	github.com/stretchr/testify v1.8.0
	gopkg.in/launchdarkly/go-sdk-common.v2 v2.5.1
	gopkg.in/launchdarkly/go-server-sdk.v5 v5.8.1
)

require (
	github.com/Shopify/sarama v1.36.0
	github.com/getsentry/sentry-go v0.13.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rs/zerolog v1.28.0
	goa.design/goa/v3 v3.8.5
)

require (
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.3 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.0.0-20220809184613-07c6da5e1ced // indirect
	golang.org/x/sys v0.0.0-20220803195053-6e608f9ce704 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/ghodss/yaml.v1 v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/aws/aws-lambda-go v1.28.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.3.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/launchdarkly/ccache v1.1.0 // indirect
	github.com/launchdarkly/eventsource v1.7.0 // indirect
	github.com/launchdarkly/go-semver v1.0.2 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	gopkg.in/launchdarkly/go-jsonstream.v1 v1.0.1 // indirect
	gopkg.in/launchdarkly/go-sdk-events.v1 v1.1.1 // indirect
	gopkg.in/launchdarkly/go-server-sdk-evaluation.v1 v1.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// These are for CVEs in these frameworks (which we don't use) and are bought in by Sentry
exclude (
	github.com/kataras/iris/v12 v12.1.8
	github.com/labstack/echo/v4 v4.1.11
)
