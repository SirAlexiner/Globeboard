// Package _func provides developer-made utility functions for use within the application.
package _func

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	authenticate "globeboard/auth"
	"globeboard/db"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Webhooks"
	"globeboard/internal/utils/structs"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	isocode = "" // Variable to store ISO code.
	title   = "" // Variable to store the title for webhook messages.
	color   = 0  // Variable to store the color code for webhook messages.
	method  = "" // Variable to store the HTTP method by which the webhook triggered.
)

// LoopSendWebhooksRegistrations sends notifications to registered webhooks about registration events.
func LoopSendWebhooksRegistrations(caller string, ci *structs.CountryInfoExternal, endpoint, eventAction string) {
	ctx := context.Background()

	// Retrieve user information; ignore error as the user is pre-validated.
	user, _ := authenticate.Client.GetUser(ctx, caller)
	email := user.DisplayName + " (" + strings.ToLower(user.Email) + ")"

	// Select appropriate message components based on the event type.
	switch eventAction {
	case Webhooks.EventRegister:
		title = Webhooks.POSTTitle
		color = Webhooks.POSTColor
		method = http.MethodPost
	case Webhooks.EventChange:
		title = Webhooks.PATCHTitle
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

	// Get the isocode from the payload data.
	isocode = ci.IsoCode

	// Fetch all webhooks from the database.
	webhooks, err := db.GetAllWebhooks()
	if err != nil {
		log.Printf("Error retrieving webhooks from database: %v", err)
		return
	}

	// Iterate through each webhook and send notifications if conditions are met.
	for _, webhook := range webhooks {
		if isRegistrationWebhookValid(caller, ci, eventAction, webhook) {
			if strings.Contains(webhook.URL, "https://discord.com") {
				sendDiscordWebhookPayload(email, title, color, method, endpoint, ci, webhook.URL)
			} else {
				sendWebhookPayload(email, title, method, endpoint, isocode, webhook.URL)
			}
		}
	}
}

// LoopSendWebhooksDashboard sends notifications to registered webhooks about dashboard events.
func LoopSendWebhooksDashboard(caller string, dr *structs.DashboardResponse) {
	ctx := context.Background()

	// Retrieve user information; ignore error as the user is pre-validated.
	user, _ := authenticate.Client.GetUser(ctx, caller)
	email := user.DisplayName + " (" + strings.ToLower(user.Email) + ")"

	// Default to INVOKE title as Dashboard endpoint GET populated dashboards at this time.
	title = Webhooks.GETTitle
	color = Webhooks.GETColor
	method = Webhooks.EventInvoke

	// Get the isocode from the payload data.
	isocode = dr.IsoCode

	// Fetch all webhooks from the database.
	webhooks, err := db.GetAllWebhooks()
	if err != nil {
		log.Printf("Error retrieving webhooks from database: %v", err)
		return
	}

	// Iterate through each webhook and send notifications if conditions are met.
	for _, webhook := range webhooks {
		if isDashboardWebhookValid(caller, dr, Webhooks.EventInvoke, webhook) {
			if strings.Contains(webhook.URL, "discord") {
				sendDiscordWebhookPayload(email, title, color, method, Endpoints.DashboardsID, dr, webhook.URL)
			} else {
				sendWebhookPayload(email, title, method, Endpoints.DashboardsID, isocode, webhook.URL)
			}
		}
	}
}

// isRegistrationWebhookValid checks if the webhook should trigger for the registration event.
func isRegistrationWebhookValid(caller string, ci *structs.CountryInfoExternal, eventAction string, webhook structs.WebhookInternal) bool {
	if webhook.UUID == "" || webhook.UUID == caller { // Validate that the webhook is associated with the user.
		// (Empty UUID is for developer webhooks)
		if webhook.Country == "" || webhook.Country == ci.IsoCode { // Validate webhook for country.
			// The event is about.
			return stringListContains(webhook.Event, eventAction) // Validate and return if webhook contains trigger event.
		}
	}
	return false
}

// isDashboardWebhookValid checks if the webhook should trigger for the dashboard event.
func isDashboardWebhookValid(caller string, dr *structs.DashboardResponse, eventAction string, webhook structs.WebhookInternal) bool {
	if webhook.UUID == "" || webhook.UUID == caller {
		if webhook.Country == "" || webhook.Country == dr.IsoCode {
			return stringListContains(webhook.Event, eventAction)
		}
	}
	return false
}

// stringListContains checks if a string is present in a slice of strings.
func stringListContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// sendDiscordWebhookPayload sends a structured message as a Discord webhook.
func sendDiscordWebhookPayload(email, title string, color int, event, endpoint string, requestBody interface{}, payloadUrl string) {
	requestBodyJSON, err := json.MarshalIndent(requestBody, "", "  ") // Pretty-print JSON for readability.
	if err != nil {
		log.Println("Error marshaling request body:", err)
		return
	}

	requestBodyString := string(requestBodyJSON)
	requestBodyString = fmt.Sprintf("```json\n%s\n```", requestBodyString) // Format JSON for Discord code block.

	// Create the payload for the Discord webhook.
	fields := []structs.Field{
		{Name: "Event", Value: event, Inline: true},
		{Name: "Endpoint", Value: endpoint, Inline: true},
		{Name: "Country", Value: isocode, Inline: true},
		{Name: "Payload", Value: requestBodyString, Inline: false},
	}

	payload := structs.WebhookPayload{
		Username:  "GlobeBoard",
		AvatarURL: "https://i.imgur.com/vjsvcxU.png",
		Embeds: []structs.Embed{
			{
				Title:       title,
				Author:      structs.Author{Name: "User: " + email},
				Description: "-------------------------------------------------------------------------------------",
				Timestamp:   time.Now().Format(time.RFC3339),
				Color:       color,
				Fields:      fields,
				Footer:      structs.Footer{Text: "Webhook Triggered:"},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload) // Serialize the payload into JSON.
	if err != nil {
		log.Println("Error marshaling payload:", err)
		return
	}

	req, err := http.NewRequest("POST", payloadUrl, bytes.NewBuffer(payloadBytes)) // Create a POST request with the payload.
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json") // Set content type to JSON.

	client := &http.Client{}
	_, err = client.Do(req) // Send the request.
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
}

// sendWebhookPayload sends a JSON formatted message to a generic webhook.
func sendWebhookPayload(email, title string, event, endpoint, country string, payloadUrl string) {
	// Create the generic webhook payload.
	payload := map[string]interface{}{
		"User":      email,
		"title":     title,
		"event":     event,
		"endpoint":  endpoint,
		"country":   country,
		"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), // Format current time in ISO8601 format.
	}

	payloadBytes, err := json.Marshal(payload) // Serialize the payload into JSON.
	if err != nil {
		log.Println("Error marshaling payload:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, payloadUrl, bytes.NewBuffer(payloadBytes)) // Create a POST request with the payload.
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json") // Set content type to JSON.

	client := &http.Client{}
	_, err = client.Do(req) // Send the request.
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
}
