package _func

import (
	"bytes"
	"encoding/json"
	"fmt"
	"globeboard/db"
	"globeboard/internal/utils/constants/Webhooks"
	"globeboard/internal/utils/structs"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func LoopSendWebhooks(ci *structs.CountryInfoPost, endpoint, eventAction string) {
	webhooks, err := db.GetWebhooks()
	if err != nil {
		log.Printf("Error retriving webhooks from database: %v", err)
		return
	}

	for _, webhook := range webhooks {
		if webhook.Country == "" || webhook.Country == ci.IsoCode {
			if stringListContains(webhook.Event, eventAction) {
				if strings.Contains(webhook.URL, "discord") {
					sendDiscordWebhookPayload(
						Webhooks.POSTTitle,
						Webhooks.POSTColor,
						http.MethodPost,
						endpoint,
						ci.IsoCode,
						ci,
						webhook.URL)
				} else {
					sendWebhookPayload(
						Webhooks.POSTTitle,
						http.MethodPost,
						endpoint,
						ci.IsoCode,
						webhook.URL)
				}
			}
		}
	}
}

func stringListContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func sendDiscordWebhookPayload(title string, color int, event, endpoint, country string, requestBody interface{}, payloadUrl string) {
	// Serialize the requestBody to a JSON string with pretty printing
	requestBodyJSON, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
		return
	}
	requestBodyString := fmt.Sprintf("```json\n%s\n```", requestBodyJSON)

	// Define default and dynamic fields
	fields := []structs.Field{
		{
			Name:   "Event",
			Value:  event,
			Inline: false,
		},
		{
			Name:   "Endpoint",
			Value:  endpoint,
			Inline: true,
		},
		{
			Name:   "Country",
			Value:  country,
			Inline: true,
		},
		{
			Name:   "Payload",
			Value:  requestBodyString,
			Inline: false,
		},
	}

	payload := structs.WebhookPayload{
		Username:  "GlobeBoard",
		AvatarURL: "https://i.imgur.com/vjsvcxU.png",
		Embeds: []structs.Embed{
			{
				Title:       title,
				Description: "-------------------------------------------------------------------------------------",
				Timestamp:   time.Now().Format(time.RFC3339), // Formatting the current time to RFC 3339
				Color:       color,
				Fields:      fields,
				Footer: structs.Footer{
					Text: "Webhook Triggered:",
				},
			},
		},
	}

	// Convert the payload into a JSON string
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		return
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", payloadUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// You can now log the response status and body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func sendWebhookPayload(title string, event, endpoint, country string, payloadUrl string) {

	payload := map[string]interface{}{
		"title":     title,
		"event":     event,
		"endpoint":  endpoint,
		"country":   country,
		"timestamp": time.Now().Format(time.RFC3339), // Formatting the current time to RFC 3339
	}

	// Convert the payload into a JSON string
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		return
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPost, payloadUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// You can now log the response status and body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
