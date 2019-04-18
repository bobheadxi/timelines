# Contributing

[![Build Status](https://dev.azure.com/bobheadxi/bobheadxi/_apis/build/status/bobheadxi.timelines?branchName=master)](https://dev.azure.com/bobheadxi/bobheadxi/_build/latest?definitionId=5&branchName=master)
[![codecov](https://codecov.io/gh/bobheadxi/timelines/branch/master/graph/badge.svg?token=8ZR61AFnLu)](https://codecov.io/gh/bobheadxi/timelines)
[![slack](https://img.shields.io/badge/slack-grey.svg?logo=slack)](https://join.slack.com/t/timelines-app/shared_invite/enQtNjEzMDE1NDk5NjAwLWZlN2ViZTE0NTNlNDZjZTNlOTNiNzZhZTZmNzgzZGVmNzcwZGE2NGJiN2QwNDQ0NzIyNmJlM2QzOTE4ZjQ3ZGE)

* [Development](#development)
  * [Commits](#commits)
  * [Web](#web)
  * [Backend](#backend)

## Development

### Commits

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

### Web

[![Netlify Status](https://api.netlify.com/api/v1/badges/b56788d9-0743-4b39-a307-66e2c99bd428/deploy-status)](https://app.netlify.com/sites/timelines-bobheadxi/deploys)

Development documentation for the Timelines web app can be found in
[`./web/README.md`](./web/README.md)

### Backend

[![GoDoc](https://godoc.org/github.com/bobheadxi/timelines?status.svg)](https://godoc.org/github.com/bobheadxi/timelines)

The backend is currently written in [Golang](https://golang.org/) and is made up
of 2 components, a server (API) and a worker. They can both be run from the
same binary, `timelines`, that you can build using:

```
make
./timelines --help
```
