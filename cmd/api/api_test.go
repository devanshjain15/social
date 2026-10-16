package main

import (
	"net/http"
	"net/http/httptest"
	"social/internal/ratelimiter"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	config := config{
		addr: ":8000",
		ratelimiter: ratelimiter.Config{
			RequestsPerTimeFrame: 20,
			TimeFrame:            time.Second * 5,
			Enabled:              true,
		},
	}
	app := NewTestApplication()
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 5

	for i := 0; i < config.ratelimiter.RequestsPerTimeFrame+marginOfError; i++ {
		request, err := http.NewRequest(http.MethodGet, ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		request.Header.Set("X-Forwarded-For", mockIP)
		response, err := client.Do(request)
		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}
		defer response.Body.Close()

		if i < config.ratelimiter.RequestsPerTimeFrame {
			if response.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", response.Status)
			}
		} else {
			if response.StatusCode != http.StatusTooManyRequests {
				t.Errorf("expected status Too Many Requests; got %v", response.Status)
			}
		}
	}

}
