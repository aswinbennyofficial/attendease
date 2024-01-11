package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aswinbennyofficial/attendease/internal/config"
	"github.com/aswinbennyofficial/attendease/internal/utility"
)

func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse and validate JWT from request
		claims, err := utility.ParseAndValidateJWT(r)
		if err != nil {
			log.Println("ERROR WHILE PARSING/VALIDATING JWT: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if the token is nearing expiration (e.g., within 2 minutes)
		expirationThreshold := time.Now().Add(2 * time.Minute)
		if claims.ExpiresAt.Time.Before(expirationThreshold) {
			// Token is nearing expiration, generate a new token
			signedToken, err := utility.GenerateToken(claims.Username)
			if err != nil {
				log.Println("ERROR OCCURRED WHILE CREATING JWT TOKEN: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Setting JWT claims for the new token
			expirationTime := time.Now().Add(time.Duration(config.LoadJwtExpiresIn()) * time.Minute)

			http.SetCookie(w, &http.Cookie{
				Name:    "JWtoken",
				Path:    "/",
				Value:   signedToken,
				Expires: expirationTime,
			})

			log.Println("TOKEN REFRESH SUCCESSFUL")
		}

		// Add claims to the request context for use in subsequent handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
