# hedera-sdk-go

### Rebuilding the Protobufs
`protoc --go_out=plugins=grpc:./hedera_proto -I ./vendor/hedera-protobuf/src/main/proto ./vendor/hedera-protobuf/src/main/proto/*.proto`
