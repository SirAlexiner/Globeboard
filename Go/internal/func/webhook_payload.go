package _func

import (
	"encoding/json"
	"fmt"
	"time"
)

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Footer struct {
	Text string `json:"text"`
}

type Embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Timestamp   string  `json:"timestamp"` // Changed to string to hold the formatted timestamp
	Color       int     `json:"color"`
	Fields      []Field `json:"fields"`
	Footer      Footer  `json:"footer"`
}

type WebhookPayload struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []Embed `json:"embeds"`
}

func MakeWebhookPayload(title string, color int, event, endpoint, country string, requestBody interface{}) WebhookPayload {
	// Serialize the requestBody to a JSON string
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
	}
	requestBodyString := fmt.Sprintf("```json\n%s\n```", requestBodyJSON)

	// Define default and dynamic fields
	fields := []Field{
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

	payload := WebhookPayload{
		Username:  "GlobeBoard",
		AvatarURL: "https://i.imgur.com/vjsvcxU.png",
		Embeds: []Embed{
			{
				Title:       title,
				Description: "-------------------------------------------------------------------------------------",
				Timestamp:   time.Now().Format(time.RFC3339), // Formatting the current time to RFC 3339
				Color:       color,
				Fields:      fields,
				Footer: Footer{
					Text: "Webhook Triggered:",
				},
			},
		},
	}
	return payload
}
