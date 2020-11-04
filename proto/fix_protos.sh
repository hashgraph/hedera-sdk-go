#!/bin/bash

# Script must be run from within the proto directory

sd 'import "(.*)";' 'import "proto/$1";' *.proto
sd '(option java_package = "com.hederahashgraph.service.proto.java";)' '$1\noption go_package = "github.com/hashgraph/hedera-sdk-go/proto";' *.proto
sd '(option java_package = "com.hederahashgraph.api.proto.java";)' '$1\noption go_package = "github.com/hashgraph/hedera-sdk-go/proto";' *.proto
sd 'google/protobuf/' '' *.proto
sd 'google.protobuf.' '' *.proto
