package xapp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHealthReadyCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/ric/v1/health/ready", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readyHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	//expected := `{"ready": true}`
	//if rr.Body.String() != expected {
	//  t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	//}
}

func TestGetHealthAliveCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/ric/v1/health/alive", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aliveHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	//expected := `{"alive": true}`
	//if rr.Body.String() != expected {
	//  t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	//}
}
