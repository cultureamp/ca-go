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

- `cipher`: easy access to kms Encrypt/Decrpyt. See [cipher](cipher/CIPHER.md) for futher details.
- `env`: easy access to common environment settings. See [env](env/ENV.md) for futher details.
- `jwt`: encode and decode the Culture Amp authentication payload. See [jwt](jwt/JWT.md) for further details.
- `log`: easy and simple logging that confirms to the logging engineer standard. See [logger](log/LOGGER.md) for further details.
- `ref`: simple methods to create pointers from literals

### Experiemental packages

These packages are under development and are subject to change.

- `x/encryption`:
- `x/kafka`:
- `x/lamdafunction`:
- `x/launchdarkly/flags`: eases the implementation and usage of LaunchDarkly for feature flags, encapsulating usage patterns in Culture Amp
- `x/log`:
- `x/request`: encapsulates the availability of request information on the request context
- `x/sarama`:
- `x/secrets`: provides methods for fetching secrets from AWS secret manager
- `x/sentry/errorreport`: eases the implementation and usage of Sentry for error reporting
- `x/valut`:

## Context

This library is the start of a replacement for
[Glamplify](https://github.com/cultureamp/glamplify). It was easier to start a
new repository and gradually move common patterns across rather than deal with a
glamplify "v2" branch, as the approach differs significantly. Keeping Glamplify around
makes it easier to migrate packages than a v2 would.

We have mindfully taken the approach of a single library with packages covering
multiple areas. This reduces maintenance, and fits the expected pattern that
most implementing services will use a reasonable proportion of the provided
functionality (given its purpose).


## Contributing

To work on `ca-go`, you'll need a working Go installation. The project currently targets Go 1.22.

### Setting up your environment

You can use [VSCode Remote Containers](https://code.visualstudio.com/docs/remote/containers) to get
up-and-running quickly. A basic configuration is defined in the `.devcontainer/`
directory. This works locally and via [GitHub Codespaces](https://github.com/features/codespaces).

#### Locally

1. Clone `ca-go` and open the directory in VSCode.
2. A prompt should appear on the bottom-right of the editor, offering to start a Remote Containers session. Click **Reopen in Container**.
3. If a prompt didn't appear, open the Command Palette (i.e. Cmd + Shift + P) and select **Remote-Containers: Open Folder in Container...**

#### Codespaces

1. Click the **Code** button above the source listing on the repository homepage.
2. Click **New codespace**.

### Pre-Commit

This is optional but is recommended for engineers working with highly sensitive secrets or data.

Config for what is checked is stored in `.pre-commit-config.yaml`

Download:
1. brew install pre-commit
2. brew install trufflehog
3. brew install snyk-cli

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
