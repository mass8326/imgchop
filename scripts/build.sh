#!/usr/bin/env bash

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
cd "$script_dir/.." || exit 1

version="0.0.0-dev"

echo "Building linux executable"
GOOS="linux" GOARCH="amd64" go build -ldflags "-X github.com/mass8326/imgchop/internal/cmd.version=$version" -o ".temp/imgchop"

echo "Building windows executable"
GOOS="windows" GOARCH="amd64" go build -ldflags "-X github.com/mass8326/imgchop/internal/cmd.version=$version" -o ".temp/imgchop.exe"
