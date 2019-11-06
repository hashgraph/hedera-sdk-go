# hedera-sdk-go

### Rebuilding the Protobufs
`protoc --go_out=plugins=grpc:./hedera_proto -I ./vendor/src/main/proto ./vendor/src/main/proto/*.proto`
