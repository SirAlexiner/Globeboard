// Package structs define structures used within the application.
package structs

type APIKey struct {
	APIKey string `json:"api_key"`
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
