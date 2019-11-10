#!/bin/sh

set -euC

# Binary name
name="zsrv"

# If you want to build for musl systems like Alpine linux
export CC=musl-gcc

go generate
go build -trimpath -ldflags "-X main.version=$(git log -n1 --format='%h_%cI')" -o "$name"

which strip >/dev/null 2>&1 \
	&& strip "$name" \
	|| echo >&2 "strip not found; skipping removal of debug symbols"

which upx >/dev/null 2>&1 \
	&& upx -qqq "$name" \
	|| echo >&2 "upx not found; skipping compression"

exit 0
