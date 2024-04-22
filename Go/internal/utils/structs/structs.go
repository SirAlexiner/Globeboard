// Package structs define data structures used within the application.
package structs

import "time"

// APIKey represents a structure for storing an API key associated with a unique identifier (UUID).
type APIKey struct {
	UUID   string `json:"uuid"`    // The unique identifier for the API key
	APIKey string `json:"api_key"` // The actual API key string
}

// CountryInfoExternal is a structure to store external-facing country information.
type CountryInfoExternal struct {
	ID         string    `json:"id"`         // Unique identifier for the country information
	Country    string    `json:"country"`    // Name of the country
	IsoCode    string    `json:"isoCode"`    // ISO code for the country
	Features   Features  `json:"features"`   // Features available for the country
	Lastchange time.Time `json:"lastchange"` // The last time the information was updated
}

// CountryInfoInternal is a structure to store internal country information.
type CountryInfoInternal struct {
	ID         string    `json:"id"`         // Unique identifier for the country information
	UUID       string    `json:"uuid"`       // An additional UUID for internal use
	Country    string    `json:"country"`    // Name of the country
	IsoCode    string    `json:"isoCode"`    // ISO code for the country
	Features   Features  `json:"features"`   // Features available for the country
	Lastchange time.Time `json:"lastchange"` // The last time the information was updated
}

// Features struct encapsulates different geographical and demographic features of a country.
type Features struct {
	Temperature      bool     `json:"temperature"`      // Boolean flag indicating retrival of temperature data
	Precipitation    bool     `json:"precipitation"`    // Boolean flag indicating retrival of precipitation data
	Capital          bool     `json:"capital"`          // Boolean flag indicating retrival of capital information
	Coordinates      bool     `json:"coordinates"`      // Boolean flag indicating retrival of geographical coordinates
	Population       bool     `json:"population"`       // Boolean flag indicating retrival of population data
	Area             bool     `json:"area"`             // Boolean flag indicating retrival of area data
	TargetCurrencies []string `json:"targetCurrencies"` // List of target currencies
}

// DashboardResponse defines the structure for dashboard service responses.
type DashboardResponse struct {
	ID            string            `json:"id"`            // Unique identifier for the dashboard entry
	Country       string            `json:"country"`       // Country name
	IsoCode       string            `json:"iso_code"`      // ISO code for the country
	Features      FeaturesDashboard `json:"features"`      // Detailed features used in the dashboard
	LastRetrieval string            `json:"lastRetrieval"` // Last retrieval time of the data
}

// FeaturesDashboard defines detailed features available on the dashboard for a country.
type FeaturesDashboard struct {
	Temperature      string                `json:"temperature,omitempty"`      // Temperature information
	Precipitation    string                `json:"precipitation,omitempty"`    // Precipitation information
	Capital          string                `json:"capital,omitempty"`          // Capital city
	Coordinates      *CoordinatesDashboard `json:"coordinates,omitempty"`      // Geographical coordinates
	Population       int                   `json:"population,omitempty"`       // Population number
	Area             string                `json:"area,omitempty"`             // Area in square kilometers
	TargetCurrencies map[string]float64    `json:"targetCurrencies,omitempty"` // Currency exchange rates
}

// CoordinatesDashboard defines latitude and longitude for a geographical location.
type CoordinatesDashboard struct {
	Latitude  string `json:"latitude,omitempty"`  // Latitude
	Longitude string `json:"longitude,omitempty"` // Longitude
}

// StatusResponse defines the current status of various APIs and services used by the application.
type StatusResponse struct {
	CountriesApi    string `json:"countries_api"` // Status of the countries API
	MeteoApi        string `json:"meteo_api"`     // Status of the meteorological API
	CurrencyApi     string `json:"currency_api"`  // Status of the currency exchange API
	FirebaseDB      string `json:"firebase_db"`   // Status of the Firebase database
	Webhooks        int    `json:"webhooks"`      // Number of active webhooks
	Version         string `json:"version"`       // Current version of the application
	UptimeInSeconds string `json:"uptime"`        // Uptime in seconds
}

// WebhookResponse defines the structure for external webhook responses.
type WebhookResponse struct {
	ID      string   `json:"id"`                // Unique identifier for the webhook
	URL     string   `json:"url"`               // URL of the webhook
	Country string   `json:"country,omitempty"` // Country associated with the webhook
	Event   []string `json:"event"`             // Events that trigger the webhook
}

// WebhookInternal is a structure to store internal-facing webhook information.
type WebhookInternal struct {
	ID      string   `json:"id"`                // Unique identifier for the webhook
	UUID    string   `json:"uuid"`              // UUID associated with the webhook
	URL     string   `json:"url"`               // URL of the webhook
	Country string   `json:"country,omitempty"` // Country associated with the webhook
	Event   []string `json:"event"`             // Events that trigger the webhook
}

// Author defines an author element for use in structured messages.
type Author struct {
	Name string `json:"name"` // Name of the author
}

// Field defines a single field in a structured message.
type Field struct {
	Name   string `json:"name"`   // Name of the field
	Value  string `json:"value"`  // Value of the field
	Inline bool   `json:"inline"` // Whether the field is inline
}

// Footer defines a footer element for use in structured messages.
type Footer struct {
	Text string `json:"text"` // Text content of the footer
}

// Embed defines the structure for embedded content in messages.
type Embed struct {
	Title       string  `json:"title"`       // Title of the embed
	Author      Author  `json:"author"`      // Author information
	Description string  `json:"description"` // Description text
	Timestamp   string  `json:"timestamp"`   // Timestamp for the embed
	Color       int     `json:"color"`       // Color code for the embed bar
	Fields      []Field `json:"fields"`      // Fields included in the embed
	Footer      Footer  `json:"footer"`      // Footer information
}

// WebhookPayload defines the payload structure for webhook messages.
type WebhookPayload struct {
	Username  string  `json:"username"`   // Username for the webhook message
	AvatarURL string  `json:"avatar_url"` // URL to the avatar image
	Embeds    []Embed `json:"embeds"`     // Embedded content in the message
}
