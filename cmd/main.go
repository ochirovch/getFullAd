package main

import (
	"aivito/getFullAd"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// helloPubSubHTTP is an HTTP handler that processes Pub/Sub messages.
func helloPubSubHTTP(w http.ResponseWriter, r *http.Request) {
	var m PubSubMessage
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode message: %v", err), http.StatusBadRequest)
		log.Printf("Could not decode message: %v", err)
		return
	}

	messageData := string(m.Message.Data)

	// Call your existing helloPubSub logic (from getPhone package)
	if err := getFullAd.HelloPubSub(context.Background(), messageData); err != nil {
		log.Printf("HelloPubSub error: %v", err)
		return
	}

	log.Printf("Successfully processed message: %s", messageData)
	fmt.Fprintln(w, "OK") // Respond to Cloud Run
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	http.HandleFunc("/", helloPubSubHTTP) // Register the HTTP handler
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
