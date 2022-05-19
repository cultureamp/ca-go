# ca-go

[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/cultureamp/ca-go)
[![License](https://img.shields.io/github/license/cultureamp/ca-go)](https://github.com/cultureamp/ca-go/blob/main/LICENSE.txt)
![Build](https://github.com/cultureamp/ca-go/workflows/pipeline/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/cultureamp/ca-go/badge.svg?branch=main)](https://coveralls.io/github/cultureamp/ca-go?branch=main)

A Go library with multiple packages to be shared by services.

This library is
intended to encapsulate the use of key practices in Culture Amp, and make their
adoption into services as straightforward as possible. The goal here is to be
light on hard opinions, but ensure that the most common patterns are supported
easily, with no hidden gotchas.

## Documentation

We use GoDoc for all documentation, including executable examples. See [the
published
documentation](https://pkg.go.dev/github.com/cultureamp/ca-go#section-directories)
on `pkg.go.dev` for more details on individual packages.

Current packages:

- `ref`: simple methods to create pointers from literals
- `launchdarkly/flags`: eases the implementation and usage of LaunchDarkly for feature flags, encapsulating usage patterns in Culture Amp
- `request`: encapsulates the availability of request information on the request context
- `sentry/errorreport`: eases the implementation and usage of Sentry for error reporting

## Context

This library is the start of a replacement for
[Glamplify](https://github.com/cultureamp/glamplify). It was easier to start a
new repository and gradually move common patterns across rather than deal with a
"v2" branch, as the approach differs significantly. Keeping Glamplify around
makes it easier to migrate packages than a v2 would.

Even though the opinions and usages here are unashamedly targeted at Culture
Amp, it's open source with an MIT license to allow for usage, adaptation or
contributions by others.

We have mindfully taken the approach of a single library with packages covering
multiple areas. This reduces maintenance, and fits the expected pattern that
most implementing services will use a reasonable proportion of the provided
functionality (given its purpose).

## Contributing

To work on `ca-go`, you'll need a working Go installation. The project currently
targets Go 1.18.

### Setting up your environment

You can use [VSCode Remote
Containers](https://code.visualstudio.com/docs/remote/containers) to get
up-and-running quickly. A basic configuration is defined in the `.devcontainer/`
directory. This works locally and via [GitHub
Codespaces](https://github.com/features/codespaces).

#### Locally

1. Clone `ca-go` and open the directory in VSCode.
2. A prompt should appear on the bottom-right of the editor, offering to start a
   Remote Containers session. Click **Reopen in Container**.
3. If a prompt didn't appear, open the Command Palette (i.e. Cmd + Shift + P)
   and select **Remote-Containers: Open Folder in Container...**

#### Codespaces

1. Click the **Code** button above the source listing on the repository
   homepage.
2. Click **New codespace**.

### Design principles

1. Aim to make the "right" way the easy way. It should be simple use this
   library for standard use cases, without being unnecessarily restrictive if
   other usage is necessary.
1. Document well. We're using GoDoc and Go's support for creating publically accessible documentation for OSS repos to our advantages here. This means that:
   1. Any public API surface should clearly self-document its intent and behaviour
   1. We make liberal use of testable `Example()` methods to make it easier to
      understand the correct usage and context of the APIs.
1. We release new packages in the `/x/` directory, moving them out when they're stable.
1. Once a package is out of `/x`, we use semantic versioning to make upgrades by
   library consumers straightforward.
