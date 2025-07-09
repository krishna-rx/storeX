package middlewares

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"main.go/utils"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "invalid token format dont found bearer", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		rawUserID, exists := claims["user_id"]
		rawRole, roleExist := claims["role"]
		if !exists || !roleExist {
			utils.ResponseError(w, http.StatusUnauthorized, "failed to get roles and userID from jwt token")
		}

		if rawRole == nil || rawUserID == nil {
			http.Error(w, "user_id  or rawRole claim missing", http.StatusUnauthorized)
			return
		}

		userID, ok := rawUserID.(string)
		role, roleIsOK := rawRole.(string)
		if !ok || !roleIsOK {
			http.Error(w, "user_id or role claim is not a string", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "roles", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
