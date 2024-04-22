// Package dashboard provides handlers for managing dashboard-related functionalities through HTTP endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/structs"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	APIInfoRetrivalError   = "Error getting country information"       // Error message for when country information cannot be retrieved.
	APICoordsRetrivalError = "Error getting Coordinates Information: " // Error message for when coordinates information cannot be retrieved.
)

// DashboardsIdHandler handles requests to the dashboard endpoint.
func DashboardsIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: // Handle GET request.
		handleDashboardGetRequest(w, r)
	default:
		// Log and return an error for unsupported HTTP methods
		log.Printf(constants.ClientConnectUnsupported, r.RemoteAddr, Endpoints.DashboardsID, r.Method)
		http.Error(w, "REST Method: "+r.Method+" not supported. Only supported methods for this endpoint is:\n"+http.MethodGet, http.StatusNotImplemented)
		return
	}
}

// handleDashboardGetRequest processes GET requests to retrieve dashboards by ID.
func handleDashboardGetRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")     // Retrieve ID from URL path.
	query := r.URL.Query()      // Extract the query parameters.
	token := query.Get("token") // Retrieve token from URL query parameters.
	if token == "" {            // Check if a token is provided.
		log.Printf(constants.ClientConnectNoToken, r.RemoteAddr, r.Method, Endpoints.DashboardsID)
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token) // Retrieve UUID associated with API token.
	if UUID == "" {                 // Check if UUID is retrieved.
		log.Printf(constants.ClientConnectUnauthorized, r.RemoteAddr, r.Method, Endpoints.DashboardsID)
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if ID == "" || ID == " " { // Check if the ID is valid.
		log.Printf(constants.ClientConnectNoID, r.RemoteAddr, r.Method, Endpoints.DashboardsID)
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	reg, err := db.GetSpecificRegistration(ID, UUID) // Retrieve registration by ID for user (UUID).
	if err != nil {
		log.Printf("%s: Error getting registration: %v", r.RemoteAddr, err)
		err := fmt.Sprintf("Dashboard doesn't exist: %v", err)
		http.Error(w, err, http.StatusNotFound)
		return
	}

	dr := new(structs.DashboardResponse) // Initialize new DashboardResponse struct.

	// Set retrieved values to response.
	dr.ID = reg.ID
	dr.Country = reg.Country
	dr.IsoCode = reg.IsoCode

	// Country information API integration.
	if getCountryInfo(w, reg, dr) {
		return
	}

	// Currency information API integration.
	if getCurrencyInfo(w, reg, dr) {
		return
	}

	// Weather information API integration.
	if getWeatherInfo(w, reg, dr) {
		return
	}

	// Set the LastRetrieval time and format it to ISO8601 format to mirror Firestore Server Timestamp.
	dr.LastRetrieval = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dr) // Encode the dashboard response into JSON and write to the response writer.
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_func.LoopSendWebhooksDashboard(UUID, dr) // Send notifications to webhooks.
}

// getWeatherInfo fetches weather information for a specific registration and updates the dashboard response.
func getWeatherInfo(w http.ResponseWriter, reg *structs.CountryInfoInternal, dr *structs.DashboardResponse) bool {
	if reg.Features.Temperature { // Check if the temperature feature is enabled.
		coords, err := _func.GetCoordinates(reg.IsoCode) // Get coordinates for the ISO code.
		if err != nil {
			log.Print(APICoordsRetrivalError, err)
			http.Error(w, APIInfoRetrivalError, http.StatusInternalServerError)
			return true
		}
		temp, err := _func.GetTemp(coords) // Get temperature for the coordinates.
		if err != nil {
			log.Print("Error getting Temperature Information: ", err)
			http.Error(w, APIInfoRetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Temperature = strconv.FormatFloat(temp, 'f', 1, 64) // Format temperature and set to dashboard response.
	}

	if reg.Features.Precipitation { // Check if the precipitation feature is enabled.
		coords, err := _func.GetCoordinates(reg.IsoCode) // Get coordinates for the ISO code.
		if err != nil {
			log.Print(APICoordsRetrivalError, err)
			http.Error(w, APIInfoRetrivalError, http.StatusInternalServerError)
			return true
		}
		precipitation, err := _func.GetPrecipitation(coords) // Get precipitation for the coordinates.
		if err != nil {
			log.Print("Error getting Temperature Information: ", err)
			http.Error(w, APIInfoRetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Precipitation = strconv.FormatFloat(precipitation, 'f', 2, 64) // Format precipitation and set to dashboard response.
	}
	return false
}

// getCurrencyInfo fetches currency exchange information for a specific registration and updates the dashboard response.
func getCurrencyInfo(w http.ResponseWriter, reg *structs.CountryInfoInternal, dr *structs.DashboardResponse) bool {
	if reg.Features.TargetCurrencies != nil && len(reg.Features.TargetCurrencies) > 0 { // Check if target-currencies feature is non nil and non-empty.
		exchangeRate, err := _func.GetExchangeRate(reg.IsoCode, reg.Features.TargetCurrencies) // Get exchange rates for the target currencies.
		if err != nil {
			log.Print("Error getting Exchange Rate Information: ", err)
			http.Error(w, APIInfoRetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.TargetCurrencies = exchangeRate // Set exchange rates to dashboard response.
	}
	return false
}

// getCountryInfo fetches country-specific information for a specific registration and updates the dashboard response.
func getCountryInfo(w http.ResponseWriter, reg *structs.CountryInfoInternal, dr *structs.DashboardResponse) bool {
	if reg.Features.Capital { // Check if the capital feature is enabled.
		capital, err := _func.GetCapital(reg.IsoCode) // Get capital for the ISO code.
		if err != nil {
			log.Print("Error getting Capital Information: ", err)
			http.Error(w, APIInfoRetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Capital = capital // Set capital to dashboard response.
	}

	if reg.Features.Coordinates { // Check if the coordinate feature is enabled.
		coords, err := _func.GetCoordinates(reg.IsoCode) // Get coordinates for the ISO code.
		if err != nil {
			log.Print(APICoordsRetrivalError, err)
			http.Error(w, APIInfoRetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Coordinates = &coords // Set coordinates to dashboard response.
	}

	if reg.Features.Population { // Check if the population feature is enabled.
		pop, err := _func.GetPopulation(reg.IsoCode) // Get population for the ISO code.
		if err != nil {
			log.Print("Error getting Population Information: ", err)
			http.Error(w, APIInfoRetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Population = pop // Set population to dashboard response.
	}

	if reg.Features.Area { // Check if area feature is enabled.
		area, err := _func.GetArea(reg.IsoCode) // Get area for the ISO code.
		if err != nil {
			log.Print("Error getting Area Information: ", err)
			http.Error(w, APIInfoRetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Area = strconv.FormatFloat(area, 'f', 1, 64) // Format area and set to dashboard response.
	}
	return false
}
