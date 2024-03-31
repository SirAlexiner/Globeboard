package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"globeboard/internal/func"
	"globeboard/internal/handlers"
	"globeboard/internal/handlers/endpoint/library"
	"globeboard/internal/handlers/endpoint/statistics"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/constants/Paths"
	"globeboard/internal/utils/structs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	// Get the port from the environment variable or set default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Define HTTP endpoints
	http.HandleFunc(Paths.Root, handlers.EmptyHandler)
	http.HandleFunc(Endpoints.Library, library.LibHandler)
	http.HandleFunc(Endpoints.ApiKey, library.APIKeyHandler)
	http.HandleFunc(Endpoints.BookCount, statistics.BookCountHandler)
	http.HandleFunc(Endpoints.Readership, statistics.ReadershipHandler)
	http.HandleFunc(Endpoints.Status, statistics.StatusHandler)
	http.HandleFunc(Endpoints.SupportedLanguages, statistics.SupportedLanguagesHandler)

	// Get supported languages
	languages, err := _func.GetSupportedLanguages()
	if err != nil {
		log.Println("Error while retrieving languages:", err)
		return
	}

	// Create main cache directory if it doesn't exist
	cacheDir := "internal/cache"
	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating main cache directory: %v\n", err)
		return
	}

	// Cache pages for each language concurrently
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(languages))

		for _, lang := range languages {
			go func(lang string) {
				defer wg.Done()
				cachePages(lang, cacheDir)
			}(lang.TwoLetter)
		}

		wg.Wait()

		log.Println("All pages cached successfully.")
	}()

	// Start the HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// cachePages caches pages for a given language.
func cachePages(languageCode, cacheDir string) {
	languageDir := filepath.Join(cacheDir, languageCode)
	if err := os.MkdirAll(languageDir, os.ModePerm); err != nil {
		log.Printf("Error creating directory for language %s: %v\n", languageCode, err)
		return
	}

	loopPaginationPages(languageCode, languageDir)
	_func.SetCacheStatus(languageCode, true)
	log.Printf("Language %s is cached", languageCode)
}

// loopPaginationPages loops over retrieval and caching of pages for a given language.
func loopPaginationPages(languageCode string, languageDir string) {
	page := 1
	for {
		// Construct filename
		filename := getPageFilename(languageDir, page)
		// Does the file exist and is it younger than a week old.
		if fileExistsAndRecent(filename) {
			// Is there a next page.
			if nextPageExists(filename) {
				page++
				continue
			} else {
				break
			}
		}

		// Fetch the json for the page and save it.
		if err := fetchAndSavePage(languageCode, page, filename); err != nil {
			log.Printf("Error processing page for language %s, page %d: %v\n", languageCode, page, err)
			return
		}

		// Break if there is no next page.
		if !nextPageExists(filename) {
			break
		}
		// increment page counter if there is a next page.
		page++
	}
}

// fetchAndSavePage retrieves the json for the language and page and saves it to the filename.
func fetchAndSavePage(languageCode string, page int, filename string) error {
	// Construct the API URL.
	apiUrl := buildAPIURL(languageCode, page)
	// Retrieve the API response.
	resp, err := fetchAPIResponse(apiUrl)
	if err != nil {
		log.Printf("Error fetching page for language %s, page %d: %v\n", languageCode, page, err)
		return err
	}

	// Confirm the API returned Status OK.
	if resp.StatusCode != http.StatusOK {
		log.Printf("Endpoint responded with not OK for language %s, page %d\n", languageCode, page)
		return errors.New("non-OK status code")
	}

	// Save the response json to filename.
	if err := saveResponseToFile(resp, filename); err != nil {
		// Handle errors with saving the json to file.
		handleFileSaveError(languageCode, page, filename, err)
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		log.Printf("Error closing body for page for language %s, page %d: %v\n", languageCode, page, err)
		return err
	}

	return nil
}

// getPageFilename returns the filename for a given language and page number.
func getPageFilename(languageDir string, page int) string {
	return filepath.Join(languageDir, fmt.Sprintf("page_%d.json", page))
}

// fileExistsAndRecent checks if a file exists and is not older than a week.
func fileExistsAndRecent(filename string) bool {
	stat, err := os.Stat(filename)
	return err == nil && time.Since(stat.ModTime()) < 7*24*time.Hour
}

// nextPageExists checks if there is a next page available.
func nextPageExists(filename string) bool {
	var response structs.BookCountResponse
	if err := decodeJSONFile(filename, &response); err != nil {
		log.Printf("Error decoding JSON file: %v\n", err)
		return false
	}
	return response.Next != ""
}

// buildAPIURL constructs the API URL for fetching pages for a language.
func buildAPIURL(languageCode string, page int) string {
	return fmt.Sprintf("%s/?languages=%s&page=%d", External.GutendexAPI, languageCode, page)
}

// fetchAPIResponse fetches the API response for the given URL.
func fetchAPIResponse(apiUrl string) (*http.Response, error) {
	return http.Get(apiUrl)
}

// saveResponseToFile saves the API response to a file.
func saveResponseToFile(resp *http.Response, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}

// handleFileSaveError handles errors encountered while saving API response to file.
func handleFileSaveError(languageCode string, page int, filename string, err error) {
	log.Printf("Error writing response JSON to file for language %s, page %d: %v\n", languageCode, page, err)
	if err := os.Remove(filename); err != nil {
		log.Printf("Error removing %s: %v\n", filename, err)
	}
}

// decodeJSONFile decodes JSON data from a file into a provided structure.
func decodeJSONFile(filename string, v interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()

	return json.NewDecoder(file).Decode(v)
}
