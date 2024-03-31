// Package statistics provides handlers for statistics-related endpoints.
package statistics

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	"globeboard/internal/func"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// ReadershipHandler handles requests to retrieve readership statistics for a specific language.
func ReadershipHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleReadershipGetRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Currently no methods are supported.", http.StatusNotImplemented)
		return
	}
}

// handleReadershipGetRequest handles GET requests to retrieve readership statistics for a specific language.
func handleReadershipGetRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	languageString := ""
	if len(parts) >= 5 {
		languageString = parts[4] // Language code will be at index 4
	} else {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Split language by comma, and validate the language code(s)
	languages := strings.Split(languageString, ",")
	if err := _func.IsValidLanguages(languages, false, false); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve bookcount information, there should only be one entry in the response
	bookCountInfo := retrieveBookCountInfo(languages)
	if len(bookCountInfo) > 1 {
		http.Error(w, "Language parameter returned more than one package of database information, expected only one", http.StatusInternalServerError)
		return
	}

	// Extract the number of authors and books from the BookCount information for the language
	numAuthors := bookCountInfo[0].Authors
	numBooks := bookCountInfo[0].Books

	// Get the optional limit parameter
	queryString := r.URL.Query().Get("limit")
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Please provide API Token", http.StatusBadRequest)
		return
	}
	exists, err := db.DoesAPIKeyExists(token)
	if err != nil {
		err := fmt.Sprintf("Error checking API key: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	if !exists {
		err := fmt.Sprintf("API key not accepted")
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	// Initialize limit as 0
	limit := 0
	// If limit is empty we update limit to be the maximum value of a 64-bit integer,
	// else we make sure the limit is a number and not 0.
	if queryString == "" {
		limit = math.MaxInt64
	} else {
		// Convert the string to an integer
		num, err := strconv.Atoi(queryString)
		if err != nil {
			http.Error(w, "Limit parameter must contain a number", http.StatusBadRequest)
			return
		}
		if num == 0 {
			http.Error(w, "Limit number must be larger than 0", http.StatusBadRequest)
			return
		}
		limit = num
	}

	// Get Country response for the language provided
	r, err = http.NewRequest(http.MethodGet, External.LanguageAPI+languages[0], nil)
	if err != nil {
		err := fmt.Sprintf("error creating request: %v", err)
		http.Error(w, err, http.StatusInternalServerError)
	}

	client := &http.Client{}
	defer client.CloseIdleConnections()

	res, err := client.Do(r)
	if err != nil {
		err := fmt.Errorf("error in response: %s", err.Error())
		log.Println(err)
	}

	var countryResponse []structs.CountryResponse
	if err := json.NewDecoder(res.Body).Decode(&countryResponse); err != nil {
		http.Error(w, "Failed to decode books response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// loop over the countries from the country response to construct readership information for each country.
	countries := loopCountryResponse(countryResponse, numBooks, numAuthors, limit)

	w.Header().Set("Content-Type", "application/json")

	// Encode countries as JSON and send the response
	if err := json.NewEncoder(w).Encode(countries); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// loopCountryResponse processes the country response and constructs readership statistics.
func loopCountryResponse(countryResponse []structs.CountryResponse, numBooks int, numAuthors int, limit int) []structs.Readership {
	// Sort countryResponse alphabetically by country name
	sort.Slice(countryResponse, func(i, j int) bool {
		return countryResponse[i].OfficialName < countryResponse[j].OfficialName
	})

	var countries []structs.Readership
	encountered := make(map[string]bool) // Map to store encountered ISO codes

	for _, country := range countryResponse {
		// Check if ISO code is already encountered
		if _, ok := encountered[country.ISO3166Alpha2]; ok {
			continue // Skip if already encountered
		}

		// Add ISO code to encountered map
		encountered[country.ISO3166Alpha2] = true
		// Retrieve the Population, Google Maps url, FLag PNG, and Alt text (For the flag image) for the country.
		population, googlemap, png, alt := _func.GetPopulationFlagMapForCountry(country.ISO3166Alpha2)

		// Construct readership statistics and append it to the existing statistics.
		countries = append(countries, structs.Readership{
			Country:    country.OfficialName,
			ISOCode:    country.ISO3166Alpha2,
			Books:      numBooks,
			Authors:    numAuthors,
			Readership: population,
			Map:        googlemap,
			Flag:       png,
			FlagAlt:    alt,
		})

		// Check if limit reached
		if len(countries) == limit {
			break
		}
	}

	return countries
}
