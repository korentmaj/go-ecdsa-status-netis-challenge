package api

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		const prefix = "Basic "
		if !strings.HasPrefix(auth, prefix) {
			http.Error(w, "Authorization header format must be Basic {base64}", http.StatusUnauthorized)
			return
		}

		encodedCredentials := strings.TrimPrefix(auth, prefix)
		// Replace "username:password" with actual credentials or a verification function
		expectedCredentials := "username:password"
		expectedEncodedCredentials := base64.StdEncoding.EncodeToString([]byte(expectedCredentials))

		if encodedCredentials != expectedEncodedCredentials {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
