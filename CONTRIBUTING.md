# Contributing

[![Build Status](https://dev.azure.com/bobheadxi/bobheadxi/_apis/build/status/bobheadxi.timelines?branchName=master)](https://dev.azure.com/bobheadxi/bobheadxi/_build/latest?definitionId=5&branchName=master)
[![codecov](https://codecov.io/gh/bobheadxi/timelines/branch/master/graph/badge.svg?token=8ZR61AFnLu)](https://codecov.io/gh/bobheadxi/timelines)
[![slack](https://img.shields.io/badge/slack-grey.svg?logo=slack)](https://join.slack.com/t/timelines-app/shared_invite/enQtNjEzMDE1NDk5NjAwLWZlN2ViZTE0NTNlNDZjZTNlOTNiNzZhZTZmNzgzZGVmNzcwZGE2NGJiN2QwNDQ0NzIyNmJlM2QzOTE4ZjQ3ZGE)

* [Development](#development)
  * [Housekeeping](#housekeeping)
    * [Commits](#commits)
    * [Branching](#branching)
  * [Web](#web)
  * [Backend](#backend)
    * [Code Style](#code-style)
    * [Development Environment](#development-environment)

## Development

### Housekeeping

#### Commits

Commits should include a single functional change, and the message should be
prefixed with the primary component being changed. For example:

```
db: implement improved not-found errors
docs: update web docs with UIKit and fixed links
```

Relevant issues should be mentioned in the commit message body:

```
chore: set up eslint for web

closes #40
```

#### Branching

Branch names should be prefixed with the primary component being changed as well.
For example:

```
git checkout -b log/resolver-middleware
```

### Web

Development documentation for the Timelines web app can be found in
[`./web/README.md`](./web/README.md).

### Backend

[![GoDoc](https://godoc.org/github.com/bobheadxi/timelines?status.svg)](https://godoc.org/github.com/bobheadxi/timelines)

The backend is currently written in [Golang](https://golang.org/) and is made up
of 2 components, a server (API) and a worker. They can both be run from the
same binary, `timelines`, that you can build using:

```
make
./timelines --help
```

Note that you'll want [Modules](https://github.com/golang/go/wiki/Modules)
enabled using `GO111MODULE=on` or by cloning this repository outside your `GOPATH`.
The minimum Go version is defined in [`go.mod`](./go.mod).

#### Code Style

All Go code should satisfy [`gofmt`](https://golang.org/cmd/gofmt/) and
[`golint`](https://github.com/golang/lint).

#### Development Environment

Install and run [Docker](https://www.docker.com/products/docker-desktop), then
use the provided [`Makefile`](./Makefile):

```
make devenv
make devpg  # initialize postgres database
```

To deploy `timelines` components using `devenv` assets, use the `--dev` flag:

```
./timelines server --dev
./timelines worker --dev
```

Utility functions for development are available as well under the `dev` command:

```
./timelines dev --help
```
