#!/usr/bin/env bash
set -euxo pipefail

FILES=`find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*/testdata/*"`
diff -u <(echo -n) <(gofmt -d -s $FILES)
diff -u <(echo -n) <(golint `go list ./... | grep -v /vendor/`)
