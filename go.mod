module github.com/cultureamp/ca-go

go 1.17

require github.com/stretchr/testify v1.7.0

require goa.design/goa v2.2.5+incompatible

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/getsentry/sentry-go v0.11.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

// These are for CVEs in these frameworks (which we don't use) and are bought in by Sentry
exclude (
	github.com/kataras/iris/v12 v12.1.8
	github.com/labstack/echo/v4 v4.1.11
)
