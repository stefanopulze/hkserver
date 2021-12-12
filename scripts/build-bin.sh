#!/usr/bin/env bash

GOOS="linux"
GOARCH="arm64"
#GOARCH="amd64"
VERSION="0.0.2"

env GOOS="$GOOS" GOARCH="$GOARCH" \
  go build -ldflags="-w -s -X configs.Version=$VERSION" -o ./output/hkserver .
