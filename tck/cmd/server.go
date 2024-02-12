package main

import (
	"fmt"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/jhttp"
	"github.com/hashgraph/hedera-sdk-go/tck/methods"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		"setup":                  handler.New(sdkService.Setup),
		"reset":                  handler.New(sdkService.Reset),
		"createAccount":          handler.New(accountService.CreateAccount),
		"createAccountFromAlias": handler.New(accountService.CreateAccountFromAlias),
		"getAccountInfo":         handler.New(accountService.GetAccountInfo),
		"updateAccountKey":       handler.New(accountService.UpdateAccountKey),
		"updateAccountMemo":      handler.New(accountService.UpdateAccountMemo),
		"deleteAccount":          handler.New(accountService.DeleteAccount),
		"generatePublicKey":      handler.New(keyService.GeneratePublicKey),
		"generatePrivateKey":     handler.New(keyService.GeneratePrivateKey),
	}

	bridge := jhttp.NewBridge(assigner, nil)

	// Listen and redirect to bridge
	http.HandleFunc("/", bridge.ServeHTTP)
	port := os.Getenv("TCK_PORT")
	if port == "" {
		port = "80"
	}
	fmt.Println("Server is listening on port: " + port)
	go func() {
		http.ListenAndServe(":"+port, nil)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	sig := <-signalCh
	fmt.Printf("Received signal: %v\n", sig)

	// Exit the application
	os.Exit(0)
}
