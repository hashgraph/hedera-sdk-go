package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/jhttp"
	"github.com/hashgraph/hedera-sdk-go/tck/methods"
	"github.com/joho/godotenv"
)

func main() {
	// Load dotenv
	_ = godotenv.Load()

	// Initialize the SDK service
	sdkService := new(methods.SDKService)
	accountService := new(methods.AccountService)
	accountService.SetSdkService(sdkService)

	// Create a new RPC server
	assigner := handler.Map{
		"setup":                  handler.New(sdkService.Setup),
		"reset":                  handler.New(sdkService.Reset),
		"createAccount":          handler.New(accountService.CreateAccount),
		"createAccountFromAlias": handler.New(accountService.CreateAccountFromAlias),
		"getAccountInfo":         handler.New(accountService.GetAccountInfo),
		"updateAccount":          handler.New(accountService.UpdateAccount),
		"deleteAccount":          handler.New(accountService.DeleteAccount),
		"generatePublicKey":      handler.New(methods.GeneratePublicKey),
		"generatePrivateKey":     handler.New(methods.GeneratePrivateKey),
	}

	bridge := jhttp.NewBridge(assigner, nil)

	// Listen and redirect to bridge
	http.HandleFunc("/", bridge.ServeHTTP)
	port := os.Getenv("TCK_PORT")
	if port == "" {
		port = "80"
	}
	fmt.Println("Server is listening on port: " + port)

	server := &http.Server{Addr: ":" + port}

	// Start the server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Wait for the termination signal
	sig := <-signalCh
	fmt.Printf("Received signal: %v\n", sig)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %s\n", err)
	}

	fmt.Println("Server shutdown complete.")
}
