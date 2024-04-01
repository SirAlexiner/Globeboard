package main

import (
	"globeboard/internal/handlers"
	"globeboard/internal/handlers/endpoint/dashboard"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Paths"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var token = os.Getenv("TOKEN")

// TestLibraryGet confirms that the Root Endpoint returns Status I'm a Teapot for All Request.
func TestRoot(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.EmptyHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusTeapot {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTeapot)
	}

	req, err = http.NewRequest("POST", Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusTeapot {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTeapot)
	}

	req, err = http.NewRequest("PUT", Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusTeapot {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTeapot)
	}
}

// TestBookCountGetLanguage confirms that the Bookcount Endpoint returns Status Bas Request for Get Request without language param.
func TestBookCountGet(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Dashboards, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestBookCountGetWrongKey confirms that the Bookcount Endpoint returns Status Not Accepted for GET Method with incorrect token.
func TestBookCountGetWrongKey(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Dashboards+"?token=c35c5742", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

// TestBookCountGetLanguageNoKey confirms that the Bookcount Endpoint returns Status Bad Request for Get Request without api key.
func TestBookCountGetLanguageNoKey(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Dashboards+"?languages=no", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestBookCountGetLanguage confirms that the Bookcount Endpoint returns Status OK for Get Request with language param.
func TestBookCountGetLanguage(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Dashboards+"?token="+token+"&languages=no", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestBookCountGetLanguageWrong confirms that the Bookcount Endpoint returns Status Bad Request for Get Request with wrongful language param.
func TestBookCountGetLanguageWrong(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Dashboards+"?token="+token+"&languages=nog", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestBookCountGetLanguages confirms that the Bookcount Endpoint returns Status OK for Get Request with multiple language param.
func TestBookCountGetLanguages(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Dashboards+"?token="+token+"&languages=no,es", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestBookCountGetLanguagesWrong confirms that the Bookcount Endpoint returns Status Bad Request for Get Request with same language param.
func TestBookCountGetLanguagesWrong(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Dashboards+"?token="+token+"&languages=no,no", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestBookcountMethodNotAllowed confirms that the Bookcount Endpoint returns Status Not Implemented for Methods other than GET.
func TestBookcountMethodNotAllowed(t *testing.T) {
	// Create a request to your endpoint with a method other than GET
	req, err := http.NewRequest("POST", Endpoints.Dashboards, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest("PUT", Endpoints.Dashboards, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

// TestReadershipGet confirms that the Notifications Endpoint returns Status Bas Request for Get Request without language param.
func TestReadershipGet(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Notifications+"?token="+token+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.NotificationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestSupportedLanguagesGetWrongKey confirms that the Notifications Endpoint returns Status Not Accepted for GET Method with incorrect token.
func TestReadershipGetWrongKey(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Notifications+"?token=c35c5742", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

// TestReadershipGetLanguageNoKey confirms that the Notifications Endpoint returns Status Bad Request for Get Request without API Token.
func TestReadershipGetLanguageNoKey(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Notifications+"no", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.NotificationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestReadershipGetLanguage confirms that the Notifications Endpoint returns Status OK for Get Request with language param.
func TestReadershipGetLanguage(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Notifications+"no/?token="+token+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.NotificationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestReadershipGetWrong confirms that the Notifications Endpoint returns Status Bad Request for Get Request with wrongful language param.
func TestReadershipGetWrong(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Notifications+"nog/?token="+token+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.NotificationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestReadershipGetLanguages confirms that the Notifications Endpoint returns Status Bad Request for Get Request with multiple language param.
func TestReadershipGetLanguages(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Notifications+"no,es/?token="+token+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.NotificationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestReadershipGetLimit confirms that the Notifications Endpoint returns Status OK for Get Request with limit param.
func TestReadershipGetLimit(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Notifications+"no/?token="+token+"&limit=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.NotificationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestReadershipGetLimitWrong confirms that the Notifications Endpoint returns Status Bad Request for Get Request with wrongful limit param.
func TestReadershipGetLimitWrong(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Notifications+"no/?token="+token+"&limit=one", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.NotificationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestReadershipMethodNotAllowed confirms that the Notifications Endpoint returns Status Not Implemented for Methods other than GET.
func TestReadershipMethodNotAllowed(t *testing.T) {
	// Create a request to your endpoint with a method other than GET
	req, err := http.NewRequest("POST", Endpoints.Notifications, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.NotificationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest("PUT", Endpoints.Notifications, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

// TestStatusGetNoKey confirms that the Status Endpoint returns Status Bad Request for GET Method without API token.
func TestStatusGetNoKey(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Status, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.StatusHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestStatusGet confirms that the Status Endpoint returns Status OK for GET Method.
func TestStatusGet(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Status+"?token="+token+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.StatusHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestStatusGetWrongKey confirms that the Status Endpoint returns Status Not Accepted for GET Method with incorrect token.
func TestStatusGetWrongKey(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Status+"?token=c35c5742", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.DashboardsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

// TestStatusMethodNotAllowed confirms that the Status Endpoint returns Status Not Implemented for Methods other than GET.
func TestStatusMethodNotAllowed(t *testing.T) {
	// Create a request to your endpoint with a method other than GET
	req, err := http.NewRequest("POST", Endpoints.Status, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.StatusHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest("PUT", Endpoints.Status, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}

// TestSupportedLanguagesGetNoKey confirms that the Supported Languages Endpoint returns Status Bad Requests for GET Method without API token.
func TestSupportedLanguagesGetNoKey(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Registrations, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.RegistrationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestSupportedLanguagesGet confirms that the Supported Languages Endpoint returns Status OK for GET Method.
func TestSupportedLanguagesGet(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Registrations+"?token="+token+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.RegistrationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestSupportedLanguagesGetWrongKey confirms that the Supported Languages Endpoint returns Status Not Accepted for GET Method with incorrect token.
func TestSupportedLanguagesGetWrongKey(t *testing.T) {
	// Create a request to your endpoint with the GET method
	req, err := http.NewRequest("GET", Endpoints.Registrations+"?token=c35c5742", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.RegistrationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

// TestSupportedLanguagesMethodNotAllowed confirms that the Supported Languages Endpoint returns Status Not Implemented for Methods other than GET.
func TestSupportedLanguagesMethodNotAllowed(t *testing.T) {
	// Create a request to your endpoint with a method other than GET
	req, err := http.NewRequest("POST", Endpoints.Registrations, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dashboard.RegistrationsHandler)

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	req, err = http.NewRequest("PUT", Endpoints.Registrations, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}
}
