package main

import (
	"context"
	"log"
	"os"

	"aivito/getFullAd"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	// Register the function with the Functions Framework.

	ctx := context.Background()
	if err := funcframework.RegisterCloudEventFunctionContext(ctx, "/", getFullAd.HelloPubSub); err != nil {
		log.Fatalf("funcframework.RegisterEventFunction: %v\n", err)
	}
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
