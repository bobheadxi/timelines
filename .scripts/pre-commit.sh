#!/usr/bin/env bash
set -o pipefail

PASS=0

# Golang checks
STAGED_GO_FILES=$(git diff --cached --name-only | grep ".go$" || true)
if [[ "$STAGED_GO_FILES" != "" ]]; then
  echo ">> info: go files were staged, running go checks"
  echo ">> check: go vet"
  go vet ./...
  if [[ $? == 1 ]]; then
    PASS=1
  fi
  echo ">> check: gofmt and golint"
  for FILE in $STAGED_GO_FILES
  do
    gofmt -d "$FILE"
    if [[ $? == 1 ]]; then
      PASS=1
    fi
    golint "-set_exit_status" "$FILE"
    if [[ $? == 1 ]]; then
      PASS=1
    fi
  done
  echo ">> info: go checks ok"
fi

if [[ $PASS == 1 ]]; then
  echo ">> error: some checks failed"
  exit 1
else
  echo ">> success: everything looks good"
  exit 0
fi
