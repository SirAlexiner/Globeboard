// Package structs define structures used within the application.
package structs

import "time"

type APIKey struct {
	APIKey string `json:"api_key"`
}

// Registrations Structs

type CountryInfoGet struct {
	Id         string    `json:"id"`
	Country    string    `json:"country"`
	IsoCode    string    `json:"isoCode"`
	Features   Features  `json:"features"`
	Lastchange time.Time `json:"lastchange"`
}

type Features struct {
	Temperature      bool     `json:"temperature"`
	Precipitation    bool     `json:"precipitation"`
	Capital          bool     `json:"capital"`
	Coordinates      bool     `json:"coordinates"`
	Population       bool     `json:"population"`
	Area             bool     `json:"area"`
	TargetCurrencies []string `json:"targetCurrencies"`
}

// StatusResponse represents the status response structure.
type StatusResponse struct {
	CountriesApi    string `json:"countries_api"`
	MeteoApi        string `json:"meteo_api"`
	CurrencyApi     string `json:"currency_api"`
	FirebaseDB      string `json:"firebase_db"`
	NotificationDb  string `json:"notification_db"`
	Webhooks        int    `json:"webhooks"`
	Version         string `json:"version"`
	UptimeInSeconds string `json:"uptime"`
}

// Webhooks Structs

type WebhookPost struct {
	URL     string   `json:"url"`
	Country string   `json:"country"`
	Event   []string `json:"event"`
}

type WebhookGet struct {
	ID      string   `json:"id"`
	URL     string   `json:"url"`
	Country string   `json:"country,omitempty"`
	Event   []string `json:"event"`
}

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
	Timestamp   string  `json:"timestamp"`
	Color       int     `json:"color"`
	Fields      []Field `json:"fields"`
	Footer      Footer  `json:"footer"`
}

type WebhookPayload struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []Embed `json:"embeds"`
}
