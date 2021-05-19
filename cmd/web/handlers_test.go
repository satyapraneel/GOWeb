package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerUp(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}
	ping(rr, r)

}
