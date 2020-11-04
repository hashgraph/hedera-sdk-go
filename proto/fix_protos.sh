#!/bin/bash

# Script must be run from within the proto directory

sd 'option java_package = "com.hederahashgraph.(api|service).proto.java";' 'option go_package = "github.com/hashgraph/hedera-sdk-go/proto";\n\noption java_package = "com.hedera.hashgraph.proto";' *.proto

sd 'import "([A-Z].*)";' 'import "proto/$1";' *.proto
