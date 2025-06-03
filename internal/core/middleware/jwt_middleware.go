package middleware

import (
	"context"
	"github.com/taninchot-work/backend-challenge/internal/constant"
	"github.com/taninchot-work/backend-challenge/internal/core/util/json"
	"github.com/taninchot-work/backend-challenge/internal/core/util/jwt"
	"log"
	"net/http"
	"strings"
)

func JwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || strings.HasPrefix("Bearer ", authHeader) {
			log.Println("Missing or malformed Authorization header")
			json.ResponseWithError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			log.Println("Invalid token")
			json.ResponseWithError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claim, err := jwt.ValidateJwt(token)
		if err != nil {
			log.Println("Token validation failed:", err)
			json.ResponseWithError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, constant.CONTEXT_KEY_USER_ID, claim.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
