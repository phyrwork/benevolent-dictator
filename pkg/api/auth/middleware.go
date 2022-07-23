package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
	"time"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{
	name: "auth",
}

type UserClaims struct {
	UserID int `json:"userId"`
	jwt.StandardClaims
}

type UserAuth struct {
	UserID int
}

var signKey *rsa.PrivateKey
var verifyKey *rsa.PublicKey

func Token(userId int, expiresIn time.Duration) (string, int, error) {
	now := time.Now()
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(expiresIn).Unix(),
			IssuedAt:  now.Unix(),
		},
		UserID: userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(signKey)
	return signedToken, int(claims.StandardClaims.ExpiresAt), err
}

func Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from header.
		header := r.Header.Get("Authorization")
		var bearer string
		if parts := strings.Split(header, "Bearer "); len(parts) > 1 {
			bearer = parts[1]
		}
		// Parse token to claims.
		var claims *UserClaims
		token, err := jwt.ParseWithClaims(bearer, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})
		if err != nil {
			log.Print(err) // TODO: improve error message
		} else if token.Valid {
			claims = token.Claims.(*UserClaims)
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

func init() {
	var err error
	// TODO: Load a persistent key from
	signKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("error generating sign key: %v", err)
	}
	verifyKey = &signKey.PublicKey
}
