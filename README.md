# hedera-sdk-go

### Dependencies
Go lang version 1.13

### Building
First [build the protobufs](#Rebuilding the Protobufs)
```
go build
```

### Rebuilding the Protobufs
`protoc --go_out=plugins=grpc:./hedera_proto -I ./vendor/hedera-protobuf/src/main/proto ./vendor/hedera-protobuf/src/main/proto/*.proto`
