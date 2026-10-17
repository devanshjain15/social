package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := NewTestApplication()
	mux := app.mount()

	t.Run("GET USER without Authenication token", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/v1/users/103", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := execRequest(request, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("GET USER with Authenication token", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/v1/users/103", nil)
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb3BoZXJzb2NpYWwiLCJzdWIiOiIxMDEiLCJhdWQiOlsiZ29waGVyc29jaWFsIl0sImV4cCI6MTc5MzE5NjU5MCwibmJmIjoxNzkwNjA0NTkwLCJpYXQiOjE3OTA2MDQ1OTB9.WDYzxVBNv8EB8fBsWZ3F50VSQnuxTEqtoW-0_FIlhQs"))
		rr := execRequest(request, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
