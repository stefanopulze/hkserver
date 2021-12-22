#!/usr/bin/env bash

GOOS="linux"
GOARCH="arm64"
#GOARCH="amd64"
VERSION=$(cat ./VERSION)

env GOOS="$GOOS" GOARCH="$GOARCH" \
  go build -ldflags="-w -s -X configs.Version=$VERSION" -o ./output/hkserver .
