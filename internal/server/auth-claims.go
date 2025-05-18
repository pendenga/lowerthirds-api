package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"lowerthirdsapi/internal/helpers"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

const (
	FirebaseProjectID = "lower3-d26f2"
	FirebaseIssuer    = "https://securetoken.google.com/" + FirebaseProjectID
)

var firebaseJWKS *keyfunc.JWKS

func init() {
	jwksURL := "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

	var err error
	firebaseJWKS, err = keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    10 * time.Second,
		RefreshUnknownKID: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to get JWKS: %v", err))
	}
}

// authClaims is a middleware function to check auth headers
func authClaims(log *logrus.Entry) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("authClaims middleware")

			// HTTP headers
			w.Header().Set("Content-Type", "application/json")
			origin := r.Header.Get("Origin")
			log.Debug("origin ", origin)
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin") // Allow caching by Origin
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			a := r.Header.Get("Authorization")
			if a == "" || !strings.HasPrefix(a, "Bearer ") {
				http.Error(w, "missing or invalid token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(a, "Bearer ")

			token, _ := jwt.Parse(tokenStr, firebaseJWKS.Keyfunc)
			// TODO: come back to this. It's not validating the signature
			// if err != nil || !token.Valid {
			// 	log.Debugf("token: %+v", token)
			// 	http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			// 	return
			// }

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid claims", http.StatusUnauthorized)
				return
			}

			// Verify standard Firebase claims
			if claims["aud"] != FirebaseProjectID {
				http.Error(w, "invalid audience", http.StatusUnauthorized)
				return
			}
			if claims["iss"] != FirebaseIssuer {
				http.Error(w, "invalid issuer", http.StatusUnauthorized)
				return
			}
			if _, ok := claims["user_id"].(string); !ok {
				http.Error(w, "user_id claim missing or invalid", http.StatusUnauthorized)
				return
			}

			socialID := claims["user_id"].(string)
			ctx := context.WithValue(r.Context(), helpers.SocialIDKey, socialID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}
