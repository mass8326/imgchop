#!/usr/bin/env bash

cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null || exit 1

version="0.0.0-dev"

echo "Building linux executable"
GOOS="linux" GOARCH="amd64" go build -ldflags "-X github.com/mass8326/imgchop/cmd.version=$version" -o "build/imgchop"

echo "Building windows executable"
GOOS="windows" GOARCH="amd64" go build -ldflags "-X github.com/mass8326/imgchop/cmd.version=$version" -o "build/imgchop.exe"
