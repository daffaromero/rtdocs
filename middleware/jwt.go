package middleware

import (
	"context"
	"net/http"
	"rtdocs/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const userContextKey contextKey = "user"

var guestTokenSecret = utils.GetEnv("GUEST_TOKEN_SECRET")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			guestClaims := jwt.MapClaims{
				"user_id":  uuid.New().String(),
				"username": "guestFromMiddleware",
				"role":     "guest",
				"exp":      time.Now().Add(time.Hour * 24).Unix(),
			}
			guestToken := jwt.NewWithClaims(jwt.SigningMethodHS256, guestClaims)
			tokenStr, err := guestToken.SignedString([]byte(guestTokenSecret))
			if err != nil {
				http.Error(w, "Failed to generate guest token", http.StatusInternalServerError)
				return
			}

			// Send the token back to the client
			w.Header().Set("Authorization", "Bearer "+tokenStr)

			// Set the user context with guest account information
			ctx := context.WithValue(r.Context(), userContextKey, guestClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		tokenStr := strings.Split(authHeader, " ")[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(guestTokenSecret), nil // Replace with env variable
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, token.Claims.(jwt.MapClaims))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
