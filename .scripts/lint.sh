#!/usr/bin/env bash
set -euxo pipefail

go vet ./...
FILES=`find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*/testdata/*"`
diff -u <(echo -n) <(gofmt -d -s $FILES)
diff -u <(echo -n) <(go run golang.org/x/lint/golint `go list ./... | grep -v /vendor/`)
