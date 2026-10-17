package main

import (
	"net/http"
	"net/http/httptest"
	"social/internal/auth"
	"social/internal/env"
	"social/internal/store"
	"social/internal/store/cache"
	"testing"

	"go.uber.org/zap"
)

func NewTestApplication() *application {
	return &application{
		logger:        zap.NewNop().Sugar(),
		store:         store.NewMockStorage(),
		cache:         cache.NewMockCacheStorage(),
		authenticator: auth.NewJWTAuthentictor(env.GetString("JWT_SECRET", ""), "gophersocial", "gophersocial"),
	}
}

func execRequest(request *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, request)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
