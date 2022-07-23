package auth

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/phyrwork/benevolent-dictator/pkg/api/graph/model"
	"net/http"
	"strings"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{
	name: "auth",
}

var Issuer = []byte("github.com/phyrwork") // TODO: Is this a reasonable definition?

type UserAuth struct {
	UserID int
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from header.
		header := r.Header.Get("Authorization")
		var bearer string
		if parts := strings.Split(header, "Bearer "); len(parts) > 1 {
			bearer = parts[1]
		}
		// Parse token to claims.
		var claims *model.UserClaims
		token, err := jwt.ParseWithClaims(bearer, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return Issuer, nil
		})
		if err == nil && token.Valid {
			claims = token.Claims.(*model.UserClaims)
		}
		// Store user auth in context.
		if claims != nil {
			userAuth := UserAuth{
				UserID: claims.UserID,
			}
			ctx := context.WithValue(r.Context(), userCtxKey, &userAuth)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func ForContext(ctx context.Context) *UserAuth {
	if raw := ctx.Value(userCtxKey); raw != nil {
		return raw.(*UserAuth)
	}
	return nil
}
