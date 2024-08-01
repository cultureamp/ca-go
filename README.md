# ca-go

[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/cultureamp/ca-go)
[![License](https://img.shields.io/github/license/cultureamp/ca-go)](https://github.com/cultureamp/ca-go/blob/main/LICENSE.txt)
![Build](https://github.com/cultureamp/ca-go/workflows/pipeline/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/cultureamp/ca-go/badge.svg?branch=main)](https://coveralls.io/github/cultureamp/ca-go?branch=main)

A Go library with multiple packages to be shared by services.

This library is intended to encapsulate the use of key practices and engineering standards of Culture Amp, and make their adoption into services as straightforward as possible. The goal here is to be light on hard opinions, but ensure that the most common patterns are supported easily, with no hidden gotchas.

## Current packages

### Stable packages

These packages are stable and there use is actively encouraged.

- `cipher`: easy access to kms Encrypt/Decrpyt. See [cipher](cipher/README.md) for further details.
- `env`: easy access to common environment settings. See [env](env/README.md) for further details.
- `jwt`: encode and decode the Culture Amp authentication payload. See [jwt](jwt/README.md) for further details.
- `kafka/consumer`: provides a subscriber (blocking) and service (non-blocking) with automatic Avro schema decoding. See [kafka/consumer](kafka/consumer/README.md) for further details.
- `launchdarkly`: eases the implementation and usage of LaunchDarkly for feature flags, encapsulating usage patterns in
Culture Amp.See [launchdarkly](launchdarkly/LAUNCHDARKLY.md) for further details.
- `log`: easy and simple logging that confirms to the logging engineer standard.
See [logger](log/README.md) for further details.
- `ref`: simple methods to create pointers from literals
- `request`: encapsulates the availability of request information on the request context.
- `secrets`: provides methods for fetching secrets from AWS secret manager.
See [secrets](secrets/README.md) for further details.
- `sentry`: eases the implementation and usage of Sentry for error reporting.
See [sentry](sentry/SENTRY.md) for further details.

## Contributing

### Installation and running with Devbox

Ensure devbox is setup as per [https://cultureamp.atlassian.net/wiki/spaces/DE/pages/3342434338/Devbox+setup](https://cultureamp.atlassian.net/wiki/spaces/DE/pages/3342434338/Devbox+setup)

Run
```
devbox run setup
```

### Pre-Commit

This is optional but is recommended for engineers working with highly sensitive secrets or data.

The pre-commit config is stored in `.pre-commit-config.yaml`.

To install / turn on for a repo:
%> pre-commit install

To uninstall / turn off for a repo:
%> pre-commit uninstall

### Design principles

1. Aim to make the "right" way the easy way. It should be simple use this library for standard use cases, without being unnecessarily restrictive if other usage is necessary.
1. Document well. This means that:
   1. Any public API surface should clearly self-document its intent and behaviour
   1. We make liberal use of testable `Example()` methods to make it easier to understand the correct usage and context of the APIs.
1. Accept interfaces, return structs.

The design of each package follows the [RFC: Design of Shared Golang packages](https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3522429030/RFC+Design+of+shared+Golang+packages)

1. Package-level methods that provide a default implementation with expected behaviours.
1. Constructor methods that allow users to implement specific versions of the package's features.
1. Constructor methods provide a clean interface for mocking behaviour.
1. Packages should have a “Testable Example” for the top level package methods.
1. Packages should not depend on each other.
