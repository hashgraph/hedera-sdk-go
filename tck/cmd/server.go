package main

import (
	"fmt"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/jhttp"
	"github.com/hashgraph/hedera-sdk-go/tck/methods"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	// Load dotenv
	_ = godotenv.Load()

	// Initialize the SDK service
	sdkService := new(methods.SDKService)
	accountService := new(methods.AccountService)
	accountService.SetSdkService(sdkService)
	keyService := new(methods.KeyService)

	// Create a new RPC server
	assigner := handler.Map{
		"setup":              handler.New(sdkService.Setup),
		"reset":              handler.New(sdkService.Reset),
		"createAccount":      handler.New(accountService.CreateAccount),
		"generatePublicKey":  handler.New(keyService.GeneratePublicKey),
		"generatePrivateKey": handler.New(keyService.GeneratePrivateKey),
	}
	bridge := jhttp.NewBridge(assigner, nil)

	// Listen and redirect to bridge
	http.HandleFunc("/", bridge.ServeHTTP)
	port := os.Getenv("TCK_PORT")
	if port == "" {
		port = "80"
	}
	fmt.Println("Server is listening on port: " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
