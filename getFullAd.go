package getFullAd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// HelloPubSub consumes a CloudEvent message and extracts the Pub/Sub message.
func HelloPubSub(ctx context.Context, messageData string) error {

	if messageData == "" {
		return fmt.Errorf("empty message data")
	}
	log.Printf("Received message: %s", messageData)
	adIDs := strings.Split(messageData, ",")
	client := &http.Client{}

	for _, adID := range adIDs {
		adID = strings.TrimSpace(adID)
		if adID == "" {
			continue
		}

		// Download ad data
		adURL := fmt.Sprintf("https://gateway.chotot.com/v1/public/ad-listing/ad_id/%s", adID)
		resp, err := client.Get(adURL)
		if err != nil {
			log.Printf("failed to get ad %s: %v", adID, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("failed to get ad %s, status: %d", adID, resp.StatusCode)
			continue
		}

		adData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed to read ad data %s: %v", adID, err)
			continue
		}

		// Send ad data to target endpoint
		postReq, err := http.NewRequest("POST", "https://raov.at/ads/ct", bytes.NewBuffer(adData))
		if err != nil {
			return fmt.Errorf("failed to create POST request: %w", err)
		}
		postReq.Header.Set("Content-Type", "application/json")

		postResp, err := client.Do(postReq)
		if err != nil {
			return fmt.Errorf("failed to POST ad: %w", err)
		}
		defer postResp.Body.Close()

		if postResp.StatusCode != http.StatusOK && postResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(postResp.Body)
			return fmt.Errorf("failed to POST ad, status: %d, body: %s", postResp.StatusCode, string(body))
		}
	}

	return nil
}
