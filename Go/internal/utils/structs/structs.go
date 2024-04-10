// Package structs define structures used within the application.
package structs

import "time"

type APIKey struct {
	UUID   string `json:"uuid"`
	APIKey string `json:"api_key"`
}

// Registrations Structs

type CountryInfoGet struct {
	ID         string    `json:"id"`
	UUID       string    `json:"uuid"`
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

// Dashboard Structs

type DashboardResponse struct {
	ID            string            `json:"id"`
	Country       string            `json:"country"`
	IsoCode       string            `json:"iso_code"`
	Features      FeaturesDashboard `json:"features"`
	LastRetrieval string            `json:"lastRetrieval"`
}

type FeaturesDashboard struct {
	Temperature      string               `json:"temperature"`
	Precipitation    string               `json:"precipitation"`
	Capital          string               `json:"capital"`
	Coordinates      CoordinatesDashboard `json:"coordinates"`
	Population       int                  `json:"population"`
	Area             string               `json:"area"`
	TargetCurrencies map[string]float64   `json:"targetCurrencies"`
}

type CoordinatesDashboard struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Status structs

// StatusResponse represents the status response structure.
type StatusResponse struct {
	CountriesApi    string `json:"countries_api"`
	MeteoApi        string `json:"meteo_api"`
	CurrencyApi     string `json:"currency_api"`
	FirebaseDB      string `json:"firebase_db"`
	Webhooks        int    `json:"webhooks"`
	Version         string `json:"version"`
	UptimeInSeconds string `json:"uptime"`
}

// Webhooks Structs

type WebhookResponse struct {
	ID      string   `json:"id"`
	URL     string   `json:"url"`
	Country string   `json:"country,omitempty"`
	Event   []string `json:"event"`
}

type WebhookGet struct {
	ID      string   `json:"id"`
	UUID    string   `json:"uuid"`
	URL     string   `json:"url"`
	Country string   `json:"country,omitempty"`
	Event   []string `json:"event"`
}

type Author struct {
	Name string `json:"name"`
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
	Author      Author  `json:"author"`
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
