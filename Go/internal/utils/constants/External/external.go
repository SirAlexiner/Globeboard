// Package External provides external API endpoints used throughout the application.
package External

// CurrencyAPI represents the endpoint for the Currency API.
// OpenMeteoAPI represents the endpoint for the Open-Meteo API.
// CountriesAPI represents the endpoint for the RESTCountries API.
const (
	CurrencyAPI  = "http://129.241.150.113:9090/currency/"
	OpenMeteoAPI = "https://api.open-meteo.com/v1/forecast"
	CountriesAPI = "http://129.241.150.113:8080/v3.1/"
)
