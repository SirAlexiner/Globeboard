package main

import (
	"globeboard/internal/handlers"
	"globeboard/internal/handlers/endpoint/dashboard"
	"globeboard/internal/utils/constants/Endpoints"
	"globeboard/internal/utils/constants/Paths"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	DisplayName = "Tester Testing"
	Email       = "Tester@Testing.test"
	Password    = "TestTesting123?!"
)

var (
	token = os.Getenv("TOKEN")
	ID    = ""
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
	err := os.Setenv("GO_ENV", "test")
	if err != nil {
		panic("Unable to set GO_ENV")
	}
}

// TestRoot confirms that Root Endpoint returns 303 See Other for All Requests.
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
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}

	req, err = http.NewRequest("POST", Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}

	req, err = http.NewRequest("PUT", Paths.Root, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}
}

// TestStatusGetNoKey confirms that the Status Endpoint returns Status Bad Request for GET Method without an API token.
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
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
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
