package middleware

import (
	"context"
	"dating-apps/app/model/entity"
	"dating-apps/helper/config"
	"dating-apps/helper/message"
	"net/http"
	"strings"

	"dating-apps/app/model/base"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	// UserIDKey is the key for the user ID value in the context.
	UserIDKey contextKey = "userID"
	bearer    string     = "bearer"
)

func Authenticate(jwtConfig *config.SecurityConfig) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tokenStr, ok := ExtractTokenFromAuthHeader(r.Header.Get("Authorization"))
			if !ok {
				base.ResponseWriter(w, http.StatusUnauthorized, base.SetDefaultResponse(r.Context(), message.UnauthorizedError))
				return
			}

			claims := &entity.Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtConfig.JwtConfig.JwtSecret), nil
			})
			if err != nil || !token.Valid {
				base.ResponseWriter(w, http.StatusUnauthorized, base.SetDefaultResponse(r.Context(), message.AuthenticationFailed))
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ExtractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], bearer) {
		return "", false
	}

	return authHeaderParts[1], true
}
