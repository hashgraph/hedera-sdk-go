#!/usr/bin/env bash

# ensure that we have protoc-gen-go
# note: requires $GOPATH/bin to be in path
# go get github.com/golang/protobuf/protoc-gen-go

# $project/proto
dir="$( cd "$( dirname $( dirname "${BASH_SOURCE[0]}" ) )" >/dev/null 2>&1 && pwd )"

protoc \
   --go_out=. --go_opt=paths=source_relative \
   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   --proto_path "$dir" \
   "$dir"/proto/*.proto

protoc \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --proto_path "$dir" \
    "$dir/proto/mirror"/*.proto
