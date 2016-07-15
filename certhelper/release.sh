#!/bin/bash

VERSION="1.0.1"
NAME=certhelper
OUTPUT=../build

echo "Building $NAME version $VERSION"

mkdir -p ${OUTPUT}

build() {
  echo -n "=> $1-$2: "
  GOOS=$1 GOARCH=$2 go build -o ${OUTPUT}/$NAME -ldflags "-X main.version=$VERSION -X main.gitHash=`git rev-parse HEAD`" ./*.go
  du -h ${OUTPUT}/${NAME}
}

build "darwin" "amd64"
