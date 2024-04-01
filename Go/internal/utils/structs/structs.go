// Package structs defines structures used within the application.
package structs

import "time"

type APIKey struct {
	APIKey string `json:"api_key"`
}

// BookResponse represents a generic response structure for book data.
type BookResponse struct {
	Count        int       `json:"count"`
	Next         string    `json:"next"`
	Previous     string    `json:"previous"`
	Results      []Book    `json:"results"`
	LastModified time.Time `json:"lastModified"`
}

// Book represents the structure of a book.
type Book struct {
	Title     string   `json:"title"`
	Authors   []Author `json:"authors"`
	Languages []string `json:"languages"`
	Formats   struct {
		TextHTML  string `json:"text/html"`
		ImageJPEG string `json:"image/jpeg"`
	} `json:"formats"`
}

// LanguageInfo represents information about a language.
type LanguageInfo struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	TwoLetter string `json:"two_letter,omitempty"`
}

// Authors represents information about book authors.
type Authors struct {
	Authors []Author `json:"authors"`
}

// Author represents the structure of an author.
type Author struct {
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	DeathYear int    `json:"death_year"`
}

// BookCountResponse represents a generic response structure for book count data.
type BookCountResponse struct {
	Count   int       `json:"count"`
	Next    string    `json:"next"`
	Results []Authors `json:"results"`
}

// BookCountInfo represents information about book counts.
type BookCountInfo struct {
	Language string  `json:"language"`
	Books    int     `json:"books"`
	Authors  int     `json:"authors"`
	Fraction float64 `json:"fraction"`
	Pages    int     `json:"pages"`
}

// CountryResponse represents a response structure for country data.
type CountryResponse struct {
	ISO3166Alpha2 string `json:"ISO3166_1_Alpha_2"`
	OfficialName  string `json:"Official_Name"`
}

// Location represents information about a location.
type Location struct {
	Maps       MapsInfo `json:"maps"`
	Flags      FlagInfo `json:"flags"`
	Population int      `json:"population"`
}

// FlagInfo represents information about flags.
type FlagInfo struct {
	PNG string `json:"png"`
	ALT string `json:"alt"`
}

// MapsInfo represents information about maps.
type MapsInfo struct {
	GoogleMaps     string `json:"googleMaps"`
	OpenStreetMaps string `json:"openStreetMaps"`
}

// Readership represents information about readership.
type Readership struct {
	Country    string `json:"country"`
	ISOCode    string `json:"isocode"`
	Books      int    `json:"books"`
	Authors    int    `json:"authors"`
	Readership int    `json:"readership"`
	Map        string `json:"map"`
	Flag       string `json:"flag"`
	FlagAlt    string `json:"flagAlt"`
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
