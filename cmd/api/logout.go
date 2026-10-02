package main

import (
	"net/http"
)

// logoutHandler godoc
//
//	@Summary		Logout user
//	@Description	Clears the authentication cookie
//	@Tags			authentication
//	@Success		204		"No content"
//	@Router			/authentication/logout [post]
func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   app.config.env == "production",
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusNoContent)
}
