#!/usr/bin/env bash

cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null || exit 1

version="v0.1.0"

echo "Building linux executable"
GOOS="linux" GOARCH="amd64" go build -ldflags "-X main.version=$version" -o "build/imgchop"

echo "Building windows executable"
GOOS="windows" GOARCH="amd64" go build -ldflags "-X main.version=$version" -o "build/imgchop.exe"
