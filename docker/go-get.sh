#!/usr/bin/env bash
GO_IMAGE="golang:1.6"
echo "docker executing \"go get $@\""
echo "--------------------------"
## install packages
docker run --rm -v "$GOPATH":/usr/libs -e GOPATH=/usr/libs "$GO_IMAGE" go get "$@"