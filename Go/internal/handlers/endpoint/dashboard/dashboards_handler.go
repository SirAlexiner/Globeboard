// Package dashboard provides handlers for dashboard-related endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DashboardsHandler handles requests to the book count API endpoint.
func DashboardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleBookCountGetRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Currently no methods are supported.", http.StatusNotImplemented)
		return
	}
}

// fetchData fetches data from the Gutendex API.
func fetchData(language string, next string) (*structs.BookCountResponse, error) {
	// Check if data is available in cache
	cacheStatus := _func.GetCacheStatus(language)
	if cacheStatus {
		// If cached data is available, load it
		cachedResponse, err := loadCachedResponse(language, next)
		if err == nil {
			return cachedResponse, nil
		}
	}

	// Create a new HTTP request
	r, err := http.NewRequest(http.MethodGet, next, nil)
	if err != nil {
		err := fmt.Errorf("error creating request: %s", err.Error())
		log.Println(err)
		return nil, err
	}

	// Add language parameter to the request if provided, and it's the initial request to the Gutendex API
	if language != "" && next == External.CurrencyAPI {
		q := r.URL.Query()
		q.Add("languages", language)
		r.URL.RawQuery = q.Encode()
	}

	// Set content type header
	r.Header.Add("content-type", "application/json")

	// Create an HTTP client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		err := fmt.Errorf("error in response: %s", err.Error())
		log.Println(err)
		return nil, err
	}

	// Decode response into BookCountResponse struct
	var response structs.BookCountResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// loadCachedResponse loads cached response from file.
func loadCachedResponse(language, next string) (*structs.BookCountResponse, error) {
	// Parse the query parameters
	u, err := url.Parse(next)
	if err != nil {
		err := fmt.Errorf("error parsing URL: %v", err)
		return nil, err
	}

	// Get the page query parameter
	page := u.Query().Get("page")

	// Construct the path to the cached JSON file
	cacheDir := filepath.Join("internal/cache", language)
	filePath := filepath.Join(cacheDir, fmt.Sprintf("page_%s.json", page))

	// Read the content of the JSON file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading cached file: %v", err)
	}

	// Unmarshal the JSON content into structs.BookCountResponse
	var response structs.BookCountResponse
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling cached JSON: %v", err)
	}

	return &response, nil
}

// handleBookCountGetRequest handles GET requests to retrieve book count information.
func handleBookCountGetRequest(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	languagesString := query.Get("languages")
	token := query.Get("token")
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
	// Split languages string into a slice of languages
	languages := strings.Split(languagesString, ",")
	// Check if the provided languages are valid
	if err := _func.IsValidLanguages(languages, false, true); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve book count information for the provided languages
	bookCountInfo := retrieveBookCountInfo(languages)

	// Sort book count information by language
	sort.Slice(bookCountInfo, func(i, j int) bool {
		return bookCountInfo[i].Language < bookCountInfo[j].Language
	})

	// Encode book count information as JSON and send the response
	err = json.NewEncoder(w).Encode(bookCountInfo)
	if err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// retrieveBookCountInfo retrieves book count information for multiple languages concurrently.
func retrieveBookCountInfo(languages []string) []structs.BookCountInfo {
	var bookCountInfo []structs.BookCountInfo
	ch := make(chan structs.BookCountInfo, len(languages)) // Use a buffered channel

	// Retrieve book count information for each language concurrently
	for _, lang := range languages {
		go fetchBookCountInfo(lang, ch)
	}

	// Collect book count information from the channel
	for range languages {
		bookCountInfo = append(bookCountInfo, <-ch)
	}
	return bookCountInfo
}

// getTotalProvidedBooks calculates the fraction of provided books.
func getTotalProvidedBooks(totalBooks int) float64 {
	// Create a new HTTP request to get the total count of books from the Gutendex API
	r, err := http.NewRequest(http.MethodGet, External.CurrencyAPI, nil)
	if err != nil {
		err := fmt.Errorf("error in creating request: %s", err.Error())
		log.Println(err)
		return 0
	}

	// Set content type header
	r.Header.Add("content-type", "application/json")

	// Create an HTTP client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		err := fmt.Errorf("error in response: %s", err.Error())
		log.Println(err)
		return 0
	}

	// Decode response into BookCountResponse struct
	var response structs.BookCountResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0
	}

	// Calculate the fraction of provided books
	fraction := float64(totalBooks) / float64(response.Count)
	// Round the fraction to four decimal places
	ratio := math.Pow(10, 4)
	fraction = math.Round(fraction*ratio) / ratio

	return fraction
}

// fetchBookCountInfo fetches book count information for a specific language.
func fetchBookCountInfo(lang string, ch chan<- structs.BookCountInfo) {
	authorMap := make(map[structs.Author]struct{})
	pageCount := 0
	nextPage := ""
	// Determine the next page URL based on cache status
	if _func.GetCacheStatus(lang) {
		nextPage = "/?page=1"
	} else {
		nextPage = External.CurrencyAPI // Initial URL
	}

	// Fetch book count information for each page until there are no more pages
	response, err := fetchData(lang, nextPage)
	if err != nil {
		log.Printf("Error fetching data for %s: %s\n", lang, err.Error())
		return
	}

	for nextPage != "" {
		response, err := fetchData(lang, nextPage)
		if err != nil {
			log.Printf("Error fetching data for %s: %s\n", lang, err.Error())
			return
		}

		// Count distinct authors in the current page
		for _, book := range response.Results {
			for _, author := range book.Authors {
				authorMap[author] = struct{}{} // Using Author struct as key, author struct includes name, birth and death
			}
		}

		pageCount++
		nextPage = response.Next // Get URL for next page
	}

	// Send book count information to the channel
	ch <- structs.BookCountInfo{
		Language: lang,
		Books:    response.Count,
		Authors:  len(authorMap),
		Fraction: getTotalProvidedBooks(response.Count),
		Pages:    pageCount,
	}
}
