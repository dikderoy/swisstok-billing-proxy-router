#!/usr/bin/env bash
GO_IMAGE="golang:1.6"
echo "docker executing \"go $@\""
echo "--------------------------"
docker run --rm -v "$GOPATH":/usr/src -w /usr/src "$GO_IMAGE" go "$@"