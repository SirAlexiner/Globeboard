// Package _func provides self-defined functions for use in the API.
package _func

import (
	"encoding/json"
	"fmt"
	"globeboard/internal/utils/constants/External"
	"globeboard/internal/utils/structs"
	"golang.org/x/text/language"
	"net/http"
)

// GetSupportedLanguages fetches supported languages from an external API.
// It returns a list of LanguageInfo structs and an error if any occurred.
func GetSupportedLanguages() ([]structs.LanguageInfo, error) {
	// Construct the URL to fetch language information.
	url := fmt.Sprintf("%sall?fields=languages", External.CountriesAPI)
	var responseData []map[string]map[string]string

	// Create an HTTP request.
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err := fmt.Errorf("error creating request: %s", err.Error())
		return nil, err
	}

	// Set content type header
	r.Header.Add("content-type", "application/json")

	// Create an HTTP client to execute the request.
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue the request and handle any errors.
	res, err := client.Do(r)
	if err != nil {
		err := fmt.Errorf("error issuing request: %s", err.Error())
		return nil, err
	}

	// Decode the JSON response.
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		err := fmt.Errorf("error JSON decoding request: %s", err.Error())
		return nil, err
	}

	// Map to store unique language codes.
	uniqueLanguages := make(map[string]string)

	// Iterate over each element in the response.
	for _, element := range responseData {
		for code, name := range element["languages"] {
			// Add the language to the map if it doesn't exist.
			if _, ok := uniqueLanguages[name]; !ok {
				uniqueLanguages[name] = code
			}
		}
	}

	// Convert the map to a list of LanguageInfo structs.
	var languages []structs.LanguageInfo
	for name, code := range uniqueLanguages {
		// Convert 3-letter code to a language tag.
		langTag, err := language.ParseBase(code)
		if err != nil {
			fmt.Printf("Error parsing language code %s: %v\n", code, err)
			continue
		}

		// Get 2-letter language code from language tag.
		twoLetterCode := langTag.String()
		if len(twoLetterCode) != 2 {
			// Skip the language if we didn't get 2-letter language code from language tag.
			continue
		}

		// Create LanguageInfo struct.
		langInfo := structs.LanguageInfo{
			Code:      code,
			Name:      name,
			TwoLetter: twoLetterCode,
		}

		// Append to the list.
		languages = append(languages, langInfo)
	}
	return languages, nil
}
