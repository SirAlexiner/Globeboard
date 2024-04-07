package _func

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	authenticate "globeboard/auth"
	"globeboard/db"
	"globeboard/internal/utils/constants/Webhooks"
	"globeboard/internal/utils/structs"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func LoopSendWebhooks(caller string, ci *structs.CountryInfoGet, endpoint, eventAction string) {
	client, err := authenticate.GetFireBaseAuthClient()
	if err != nil {
		log.Printf("Error initializing Firebase Auth: %v", err)
		return
	}

	ctx := context.Background()

	// Ignoring error as we've already confirmed the caller at the endpoint.
	user, _ := client.GetUser(ctx, caller)

	email := user.Email
	title := ""
	color := 0
	method := ""

	switch eventAction {
	case Webhooks.EventRegister:
		title = Webhooks.POSTTitle
		color = Webhooks.POSTColor
		method = http.MethodPost
	case Webhooks.EventChange:
		title = Webhooks.PUTTitle
		color = Webhooks.PUTColor
		method = http.MethodPatch
	case Webhooks.EventDelete:
		title = Webhooks.DELETETitle
		color = Webhooks.DELETEColor
		method = http.MethodDelete
	case Webhooks.EventInvoke:
		title = Webhooks.GETTitle
		color = Webhooks.GETColor
		method = http.MethodGet

	}
	webhooks, err := db.GetAllWebhooks()
	if err != nil {
		log.Printf("Error retriving webhooks from database: %v", err)
		return
	}

	for _, webhook := range webhooks {
		if isWebhookValid(caller, ci, eventAction, webhook) && strings.Contains(webhook.URL, "discord") {
			sendDiscordWebhookPayload(
				email,
				title,
				color,
				method,
				endpoint,
				ci,
				webhook.URL)
		} else {
			sendWebhookPayload(
				email,
				title,
				method,
				endpoint,
				ci.IsoCode,
				webhook.URL)
		}
	}
}

func isWebhookValid(caller string, ci *structs.CountryInfoGet, eventAction string, webhook structs.WebhookGet) bool {
	if webhook.UUID == "" || webhook.UUID == caller {
		if webhook.Country == "" || webhook.Country == ci.IsoCode {
			if stringListContains(webhook.Event, eventAction) {
				return true
			}
			return false
		}
		return false
	}
	return false
}

func stringListContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func sendDiscordWebhookPayload(email, title string, color int, event, endpoint string, requestBody *structs.CountryInfoGet, payloadUrl string) {
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
			Inline: true,
		},
		{
			Name:   "Endpoint",
			Value:  endpoint,
			Inline: true,
		},
		{
			Name:   "Country",
			Value:  requestBody.IsoCode,
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
				Title: title,
				Author: structs.Author{
					Name: "User: " + email,
				},
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

func sendWebhookPayload(email, title string, event, endpoint, country string, payloadUrl string) {

	payload := map[string]interface{}{
		"email":     email,
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
	log.Println("Response Status:" + resp.Status + "Response Body:" + string(body))
}
