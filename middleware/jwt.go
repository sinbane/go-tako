package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sinbane/tako/config"
)

func JWT(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request URL is in the bypass URLs
			for _, bypassURL := range cfg.JWT.BypassURLs {
				if r.URL.Path == bypassURL {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Extract the token from the Authorization header
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			// The token should be in the format "Bearer <token>"
			tokens := strings.Split(auth, " ")
			if len(tokens) != 2 || tokens[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			// Parse and validate the token
			token, err := jwt.Parse(tokens[1], func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Return the secret key used to sign the token
				return []byte(cfg.JWT.Secret), nil
			})

			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// If the token is valid, you can extract claims and add them to the request context
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				// Add claims to request context
				ctx := context.WithValue(r.Context(), "user", claims)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
