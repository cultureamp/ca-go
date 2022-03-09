module github.com/cultureamp/ca-go

go 1.17

require (
	github.com/aws/aws-sdk-go v1.42.7
	github.com/launchdarkly/go-server-sdk-dynamodb v1.1.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/launchdarkly/go-sdk-common.v2 v2.5.0
	gopkg.in/launchdarkly/go-server-sdk.v5 v5.8.1
)

require (
	github.com/getsentry/sentry-go v0.11.0
	goa.design/goa/v3 v3.6.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/launchdarkly/ccache v1.1.0 // indirect
	github.com/launchdarkly/eventsource v1.7.0 // indirect
	github.com/launchdarkly/go-semver v1.0.2 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/launchdarkly/go-jsonstream.v1 v1.0.1 // indirect
	gopkg.in/launchdarkly/go-sdk-events.v1 v1.1.1 // indirect
	gopkg.in/launchdarkly/go-server-sdk-evaluation.v1 v1.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

// These are for CVEs in these frameworks (which we don't use) and are bought in by Sentry
exclude (
	github.com/kataras/iris/v12 v12.1.8
	github.com/labstack/echo/v4 v4.1.11
)
