// Package dashboard provides handlers for dashboard-related endpoints.
package dashboard

import (
	"encoding/json"
	"fmt"
	"globeboard/db"
	_func "globeboard/internal/func"
	"globeboard/internal/utils/structs"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	RetrivalError = "Error getting country information"
)

// DashboardsHandler handles requests to the book count API endpoint.
func DashboardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleDashboardGetRequest(w, r)
	default:
		http.Error(w, "REST Method: "+r.Method+" not supported. Currently only"+http.MethodGet+" methods are supported.", http.StatusNotImplemented)
		return
	}
}

// handleDashboardGetRequest handles GET requests to retrieve book count information.
func handleDashboardGetRequest(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("ID")
	query := r.URL.Query()
	token := query.Get("token")
	if token == "" {
		http.Error(w, ProvideAPI, http.StatusUnauthorized)
		return
	}
	UUID := db.GetAPIKeyUUID(token)
	if UUID == "" {
		err := fmt.Sprintf(APINotAccepted)
		http.Error(w, err, http.StatusNotAcceptable)
		return
	}
	if ID == "" {
		http.Error(w, ProvideID, http.StatusBadRequest)
		return
	}

	reg, err := db.GetSpecificRegistration(ID, UUID)
	if err != nil {
		err := fmt.Sprint("Document doesn't exist: ", err)
		http.Error(w, err, http.StatusNotFound)
		return
	}

	dr := new(structs.DashboardResponse)

	dr.ID = reg.ID
	dr.Country = reg.Country
	dr.IsoCode = reg.IsoCode

	// Countries API
	if getCountryInfo(w, reg, dr) {
		return
	}

	// Currency API
	if getCurrencyInfo(w, reg, dr) {
		return
	}

	// Open-Meteo API
	if getWeatherInfo(w, reg, dr) {
		return
	}

	// Set the LastRetrieval time and format it to ISO8601 format to mirror Firestore Timestamp
	dr.LastRetrieval = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response
	w.WriteHeader(http.StatusOK)

	// Serialize the struct to JSON and write it to the response
	err = json.NewEncoder(w).Encode(dr)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_func.LoopSendWebhooksDashboard(UUID, dr)
}

func getWeatherInfo(w http.ResponseWriter, reg *structs.CountryInfoGet, dr *structs.DashboardResponse) bool {
	if reg.Features.Temperature {
		temp, err := _func.GetTemp(dr.Features.Coordinates)
		if err != nil {
			log.Print("Error getting Temperature Information: ", err)
			http.Error(w, RetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Temperature = strconv.FormatFloat(temp, 'f', 1, 64)
	}

	if reg.Features.Precipitation {
		precipitation, err := _func.GetPrecipitation(dr.Features.Coordinates)
		if err != nil {
			log.Print("Error getting Temperature Information: ", err)
			http.Error(w, RetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Precipitation = strconv.FormatFloat(precipitation, 'f', 2, 64)
	}
	return false
}

func getCurrencyInfo(w http.ResponseWriter, reg *structs.CountryInfoGet, dr *structs.DashboardResponse) bool {
	if reg.Features.TargetCurrencies != nil && len(reg.Features.TargetCurrencies) != 0 {
		exchangeRate, err := _func.GetExchangeRate(reg.IsoCode, reg.Features.TargetCurrencies)
		if err != nil {
			log.Print("Error getting Exchange Rate Information: ", err)
			http.Error(w, RetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.TargetCurrencies = exchangeRate
	}
	return false
}

func getCountryInfo(w http.ResponseWriter, reg *structs.CountryInfoGet, dr *structs.DashboardResponse) bool {
	if reg.Features.Capital {
		capital, err := _func.GetCapital(reg.IsoCode)
		if err != nil {
			log.Print("Error getting Capital Information: ", err)
			http.Error(w, RetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Capital = capital
	}

	if reg.Features.Coordinates {
		coords, err := _func.GetCoordinates(reg.IsoCode)
		if err != nil {
			log.Print("Error getting Coordinates Information: ", err)
			http.Error(w, RetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Coordinates = coords
	}

	if reg.Features.Population {
		pop, err := _func.GetPopulation(reg.IsoCode)
		if err != nil {
			log.Print("Error getting Population Information: ", err)
			http.Error(w, RetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Population = pop
	}

	if reg.Features.Area {
		area, err := _func.GetArea(reg.IsoCode)
		if err != nil {
			log.Print("Error getting Area Information: ", err)
			http.Error(w, RetrivalError, http.StatusInternalServerError)
			return true
		}
		dr.Features.Area = strconv.FormatFloat(area, 'f', 1, 64)
	}
	return false
}
