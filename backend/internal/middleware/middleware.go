package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Cakra17/imphnen/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type userClaimsKey struct{}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			utils.ResponseJson(w, http.StatusUnauthorized, utils.Response{
				Message: "Token tidak ditemukan",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			utils.ResponseJson(w, http.StatusUnauthorized, utils.Response{
				Message: "Format token tidak sesuai",
			})
			return
		}

		secret := os.Getenv("JWT_SECRET")
		token, err := utils.ValidateToken(tokenStr, secret)
		if err != nil {
			log.Printf("%s", err.Error())
			utils.ResponseJson(w, http.StatusUnauthorized, utils.Response{
				Message: "Token Invalid",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.ResponseJson(w, http.StatusUnauthorized, utils.Response{
				Message: "Token tidak mennyimpan user info",
			})
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(ctx context.Context) (jwt.MapClaims, bool) {
	val := ctx.Value(userClaimsKey{})
	claims, ok := val.(jwt.MapClaims)
	return claims, ok
}
