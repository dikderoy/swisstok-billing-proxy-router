#!/usr/bin/env bash
GO_IMAGE="golang:1.6"
APP_NAME="HAR"
VERSION=`cat .semver`
rm ./build/*
docker run --rm \
	-v "$GOPATH":/usr/gopath \
	-e GOPATH=/usr/gopath \
	-v $PWD:/usr/src/app \
	-w /usr/src/app \
	"$GO_IMAGE" \
	go build -v \
	-o ./build/"$APP_NAME-$VERSION";