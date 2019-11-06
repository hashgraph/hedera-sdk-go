# hedera-sdk-go

### Rebuilding the Protobufs
`protoc --go_out=plugins=grpc:./hedera_proto -I ./hedera-protobuf/src/main/proto ./hedera-protobuf/src/main/proto/*.proto`
