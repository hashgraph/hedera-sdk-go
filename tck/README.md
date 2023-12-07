# Go SDK TCK server

This is a server that implements the [SDK TCK specification](https://github.com/hashgraph/hedera-sdk-tck/) for the Go SDK.

## Running the server

To run the server you need to run

```
go run cmd/server.go
```

This will start the server on port **80**. You can change the port by setting the `TCK_PORT` environment variable or by adding a .env file with the same variable.
