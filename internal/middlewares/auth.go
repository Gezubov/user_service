package middlewares

import (
	"context"
	"net/http"

	"github.com/Gezubov/user_service/config"
	"github.com/golang-jwt/jwt/v5"
)

type key string

const UserIDKey key = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "Invalid token payload", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
