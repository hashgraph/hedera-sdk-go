#!/usr/bin/env bash

# ensure that we have protoc-gen-go
# note: requires $GOPATH/bin to be in path
go get github.com/golang/protobuf/protoc-gen-go

# $project/proto
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

protoc \
  --go_out=plugins=grpc:"$dir" \
  -I "$dir"/. \
  "$dir"/*.proto

protoc \
  --go_out=plugins=grpc:"$dir"/mirror \
  --proto_path "$dir"/mirror \
  -I "$dir" \
  "$dir"/mirror/*.proto
