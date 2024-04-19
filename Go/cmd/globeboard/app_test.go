package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"globeboard/internal/handlers"
	"globeboard/internal/handlers/endpoint/dashboard"
	"globeboard/internal/handlers/endpoint/util"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Paths"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

const (
	DisplayName = "Tester Testing"
	Email       = "Tester@Testing.test"
	Password    = "TestTesting123?!"
)

var (
	mux        = http.NewServeMux()
	wrongToken = "bhuiozdfbbjkwsrbnjlsfbjnklsdv" //Keyboard Mash
	token      = "sk-token-brrr-access"
	UUID       = "me_me_me_me"
	docId1     = "420"
	docId2     = "420"
	webhookId1 = "69"
	webhookId2 = "69"
)

func fileExistsTest(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func init() {
	if !fileExistsTest(os.Getenv("FIREBASE_CREDENTIALS_FILE")) {
		log.Panic("Firebase Credentials file is not mounted: ", os.Getenv("FIREBASE_CREDENTIALS_FILE"))
	}

	mux.HandleFunc(Paths.Root, handlers.EmptyHandler)
	mux.HandleFunc(Endpoints.UserRegistration, util.UserRegistrationHandler)
	mux.HandleFunc(Endpoints.UserDeletionId, util.UserDeletionHandler)
	mux.HandleFunc(Endpoints.ApiKey, util.APIKeyHandler)
	mux.HandleFunc(Endpoints.RegistrationsID, dashboard.RegistrationsIdHandler)
	mux.HandleFunc(Endpoints.Registrations, dashboard.RegistrationsHandler)
	mux.HandleFunc(Endpoints.DashboardsID, dashboard.DashboardsIdHandler)
	mux.HandleFunc(Endpoints.NotificationsID, dashboard.NotificationsIdHandler)
	mux.HandleFunc(Endpoints.Notifications, dashboard.NotificationsHandler)
	mux.HandleFunc(Endpoints.Status, dashboard.StatusHandler)

}

/* Test Root Once */

// TestRoot confirms that Root Endpoint returns 303 See Other for All Requests.
func TestRoot(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}

	req, err = http.NewRequest(http.MethodPost, Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}

	req, err = http.NewRequest(http.MethodPut, Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}

	req, err = http.NewRequest(http.MethodPatch, Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}

	req, err = http.NewRequest(http.MethodDelete, Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}
}

/* Run tests as intended */

func TestRegisterHandlerRegister(t *testing.T) {
	form := url.Values{}
	form.Add("username", DisplayName)
	form.Add("email", Email)
	form.Add("password", Password)

	req, err := http.NewRequest(http.MethodPost, Endpoints.UserRegistration, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Decode the JSON response
	var response struct {
		Token  string `json:"token"`
		UserID string `json:"userid"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body:", err)
	}

	token = response.Token
	UUID = response.UserID
}

func TestDeleteAPIKeyHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.ApiKey+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", UUID)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("DELETE handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestGetAPIKeyHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, Endpoints.ApiKey, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", UUID)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("GET handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Optionally decode the response to check if the correct API key is retrieved
	var response struct {
		APIKey string `json:"token"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body:", err)
	}

	token = response.APIKey
}

func TestStatusGet(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Status+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestNotificationsHandlerPostDiscord(t *testing.T) {
	notificationData := []byte(`{
		"url": "https://localhost/discord",
		"country": "",
		"event": ["INVOKE","REGISTER","CHANGE","DELETE"]
	}`)
	req, err := http.NewRequest(http.MethodPost, Endpoints.Notifications+"?token="+token, bytes.NewBuffer(notificationData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body:", err)
	}

	webhookId1 = response.ID
}

func TestNotificationsHandlerPost(t *testing.T) {
	notificationData := []byte(`{
		"url": "https://localhost/",
		"country": "",
		"event": ["INVOKE","DELETE"]
	}`)
	req, err := http.NewRequest(http.MethodPost, Endpoints.Notifications+"?token="+token, bytes.NewBuffer(notificationData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body:", err)
	}

	webhookId2 = response.ID
}

func TestNotificationsHandlerIdGet(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, Endpoints.Notifications+"/"+webhookId1+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestNotificationsHandlerGet(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, Endpoints.Notifications+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRegistrationsHandlerPost(t *testing.T) {
	registrationData := []byte(`{
		"country": "united states",
		"features": { 
			"temperature": true,
			"coordinates": true
		   }
	}`)

	req, err := http.NewRequest(http.MethodPost, Endpoints.Registrations+"?token="+token, bytes.NewBuffer(registrationData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response struct {
		ID         string `json:"id"`
		LastChange string `json:"lastChange"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body:", err)
	}

	docId1 = response.ID
}

func TestRegistrationsHandlerPostMinimal(t *testing.T) {
	registrationData := []byte(`{
		"isoCode": "no",
		"features": { 
			"temperature": true
		}
	}`)

	req, err := http.NewRequest(http.MethodPost, Endpoints.Registrations+"?token="+token, bytes.NewBuffer(registrationData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response struct {
		ID         string `json:"id"`
		LastChange string `json:"lastChange"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body:", err)
	}

	docId2 = response.ID
}

func TestRegistrationsHandlerGet(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, Endpoints.Registrations+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRegistrationsIdHandlerPatch(t *testing.T) {
	patchData := []byte(`{
		"features": { 
			"temperature": true,
			"precipitation": true,
			"capital": true,
			"coordinates": true,
			"population": true,
			"area": true,
			"targetCurrencies": ["jpy", "nok", "eur","gbp"]
		}
    }`)

	req, err := http.NewRequest(http.MethodPatch, Endpoints.Registrations+"/"+docId1+"?token="+token, bytes.NewBuffer(patchData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}
}

func TestRegistrationsIdHandlerGet(t *testing.T) {
	testUrl := fmt.Sprintf("%s/%s?token=%s", Endpoints.Registrations, docId1, token)
	req, err := http.NewRequest(http.MethodGet, testUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDashboardIdHandlerGet(t *testing.T) {
	testUrl := fmt.Sprintf("%s/%s?token=%s", Endpoints.Dashboards, docId1, token)
	req, err := http.NewRequest(http.MethodGet, testUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	t.Log(rr.Body.String())
}

func TestDashboardIdHandlerGetMinimal(t *testing.T) {
	testUrl := fmt.Sprintf("%s/%s?token=%s", Endpoints.Dashboards, docId2, token)
	req, err := http.NewRequest(http.MethodGet, testUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	t.Log(rr.Body.String())
}

func TestRegistrationsIdHandlerDeleteMinimal(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Registrations+"/"+docId2+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

/* Run test with wrong token */

func TestDeleteAPIKeyHandlerWrongToken(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.ApiKey+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", UUID)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DELETE handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestStatusGetWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Status+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestRegistrationsPostWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodPost, Endpoints.Registrations+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestRegistrationsGetWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Registrations+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestRegistrationsGetIdWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Registrations+"/"+docId1+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestRegistrationsPatchIdWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodPatch, Endpoints.Registrations+"/"+docId1+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestRegistrationsDeleteIdWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Registrations+"/"+docId1+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestDashboardGetIdWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Dashboards+"/"+docId1+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestNotificationsPostWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodPost, Endpoints.Notifications+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestNotificationsGetWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Notifications+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestNotificationsGetIdWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Notifications+"/"+webhookId1+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

func TestNotificationsDeleteIdWrongToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Notifications+"/"+webhookId1+"?token="+wrongToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

/* Run tests without a token */

func TestDeleteAPIKeyHandlerNoToken(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.ApiKey, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", UUID)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("DELETE handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestStatusGetNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Status, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestRegistrationsPostNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodPost, Endpoints.Registrations, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestRegistrationsGetNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Registrations, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestRegistrationsGetIdNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Registrations+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestRegistrationsPatchIdNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodPatch, Endpoints.Registrations+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestRegistrationsDeleteIdNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Registrations+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestDashboardGetIdNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Dashboards+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestNotificationsPostNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodPost, Endpoints.Notifications, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestNotificationsGetNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Notifications, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestNotificationsGetIdNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Notifications+"/"+webhookId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestNotificationsDeleteIdNoToken(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Notifications+"/"+webhookId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

/* Empty ID */

func TestRegistrationsGetEmptyId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Registrations+"/?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestRegistrationsPatchEmptyId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodPatch, Endpoints.Registrations+"/?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestRegistrationsDeleteEmptyId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Registrations+"/?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDashboardGetEmptyId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Dashboards+"/?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestNotificationsGetEmptyId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Notifications+"/?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestNotificationsDeleteEmptyId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Notifications+"/?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

/* Wrong ID */

func TestRegistrationsGetWrongId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Registrations+"/aaaaaaaaaaaaaaaaa?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestRegistrationsPatchWrongId(t *testing.T) {
	patchData := []byte(`{
        "features": {
            "targetCurrencies": ["EUR", "USD", "NOK"]
        }
    }`)
	req, err := http.NewRequest(http.MethodPatch, Endpoints.Registrations+"/aaaaaaaaaaaaaaaaa?token="+token, bytes.NewBuffer(patchData))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestRegistrationsDeleteWrongId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Registrations+"/aaaaaaaaaaaaaaaaa?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDashboardGetWrongId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Dashboards+"/aaaaaaaaaaaaaaaaa?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestNotificationsGetWrongId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.Notifications+"/aaaaaaaaaaaaaaaaa?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestNotificationsDeleteWrongId(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Notifications+"/aaaaaaaaaaaaaaaaa?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

/* Register Wrong Email, Password */

func TestRegisterHandlerRegisterBadEmail(t *testing.T) {
	form := url.Values{}
	form.Add("username", DisplayName)
	form.Add("email", "TesterTesting.test")
	form.Add("password", Password)

	req, err := http.NewRequest(http.MethodPost, Endpoints.UserRegistration, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestRegisterHandlerRegisterBadPassword(t *testing.T) {
	form := url.Values{}
	form.Add("username", DisplayName)
	form.Add("email", Email)
	form.Add("password", "password")

	req, err := http.NewRequest(http.MethodPost, Endpoints.UserRegistration, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

/* Empty POST/PATCH Body */

func TestDeleteAPIKeyHandlerEmpty(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.ApiKey+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("DELETE handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestGetAPIKeyHandlerEmpty(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, Endpoints.ApiKey, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("GET handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestRegisterHandlerEmptyRegister(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, Endpoints.UserRegistration, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestNotificationsHandlerEmptyPost(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, Endpoints.Notifications+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestRegistrationsHandlerEmptyPost(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, Endpoints.Registrations+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestRegistrationsIdHandlerEmptyPatch(t *testing.T) {
	req, err := http.NewRequest(http.MethodPatch, Endpoints.Registrations+"/"+docId1+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

/* Delete User No UUID & Wrong UUID */

func TestRegisterHandlerDeleteNoUUID(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.UserDeletion+"/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestRegisterHandlerDeleteWrongUUID(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.UserDeletion+"/NTNU2024", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

/* Endpoint "Not Implemented" Methods Check */

func TestUserRegister(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.UserRegistration, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPut, Endpoints.UserRegistration, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPatch, Endpoints.UserRegistration, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodDelete, Endpoints.UserRegistration, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

func TestUserDeletion(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest(http.MethodGet, Endpoints.UserDeletion+"/"+UUID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPut, Endpoints.UserDeletion+"/"+UUID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPatch, Endpoints.UserDeletion+"/"+UUID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPost, Endpoints.UserDeletion+"/"+UUID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

func TestAPIKeyHandler(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPut, Endpoints.ApiKey, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPatch, Endpoints.ApiKey, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPost, Endpoints.ApiKey, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

func TestStatus(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPut, Endpoints.Status, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPatch, Endpoints.Status, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPost, Endpoints.Status, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodDelete, Endpoints.Status, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

func TestRegistrations(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPut, Endpoints.Registrations, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPatch, Endpoints.Registrations, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodDelete, Endpoints.Registrations, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

func TestRegistrationsId(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPut, Endpoints.Registrations+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPost, Endpoints.Registrations+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

func TestDashboardId(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPut, Endpoints.Dashboards+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPost, Endpoints.Dashboards+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPatch, Endpoints.Dashboards+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodDelete, Endpoints.Dashboards+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

func TestNotifications(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPut, Endpoints.Notifications, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPatch, Endpoints.Notifications, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodDelete, Endpoints.Notifications, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

func TestNotificationsId(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPut, Endpoints.Notifications+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPost, Endpoints.Notifications+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest(http.MethodPatch, Endpoints.Notifications+"/"+docId1, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	mux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

/* Clean UP */

func TestRegistrationsIdHandlerDelete(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Registrations+"/"+docId1+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestNotificationsHandlerDeleteDiscord(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Notifications+"/"+webhookId1+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestNotificationsHandlerDelete(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.Notifications+"/"+webhookId2+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestDeleteAPIKey(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.ApiKey+"?token="+token, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", UUID)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("DELETE handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestRegisterHandlerDelete(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, Endpoints.UserDeletion+"/"+UUID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
