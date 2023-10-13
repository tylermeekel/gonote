package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int
	jwt.RegisteredClaims
}

// signJWT signs a JWT using an expiry time and username as part of the
// RegisteredClaims. It then returns the signed string.
func signJWT(userID int, expTime time.Time) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// parseJWT takes a tokenString as a parameter, and checks if it is valid.
// It then returns the userID and an error.
func parseJWT(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// Verify that the signing method is HMAC-SHA256 and return the secret key.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	} else {
		return 0, errors.New("invalid token")
	}
}

func getUserIDFromContext(r *http.Request) int {
	userID, ok := r.Context().Value(userIDKey).(int)
	if userID == 0 || !ok {
		return 0
	}

	return userID
}

//Middleware

func checkAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			ctx := context.WithValue(r.Context(), userIDKey, 0)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, err := parseJWT(token.Value)
		if err != nil {
			ctx := context.WithValue(r.Context(), userIDKey, 0)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
