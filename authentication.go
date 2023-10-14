package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

func (app *App) authRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/register", app.handleRegisterUser)
	router.Post("/login", app.handleLoginUser)
	router.Post("/logout", app.handleLogoutUser)

	return router
}

// handleRegisterUser takes the username and password from the form request.
// It then hashes the password using bcrypt, and calls the createUser function
// with the username and hashed password. If everything goes well it sends a toast
// message back to the user to let them know.
func (app *App) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 10)
	if err != nil {
		fmt.Println(err.Error())
	}

	user, err := app.createUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.sendErrorToast(w, "Error: Internal Server Error")
	} else {
		fmt.Printf("Created user \"%s\"\n", user.Username)
		sendToast(w, "User successfully created")
	}

}

// handleLogin user takes the username and password from the form request.
// It then queries the database for the given username, and checks the password
// against the hash that is stored in the database. If the username exists in the database
// and the hash matches, it then generates a JWT with a 10 minute lifetime, using the username
// as part of the claims, and sends the JWT back as a cookie to the user.
func (app *App) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	queriedUser, err := app.getUserByUsername(username)
	if err != nil{
		app.sendErrorToast(w, "Internal server error")
		return
	}

	//Check hash and password and return an error if they do not match
	err = bcrypt.CompareHashAndPassword([]byte(queriedUser.Password), []byte(password))
	if err != nil || queriedUser.Username == "" {
		app.sendErrorToast(w, "Incorrect username or password")
		return
	}

	//set expiration time and create JWT using user ID and expiration time and return error if there is
	//an error generating the jwt
	expirationTime := time.Now().Add(10 * time.Minute)
	signedString, err := signJWT(queriedUser.ID, expirationTime)
	if err != nil {
		app.log.Println(err.Error())
		app.sendErrorToast(w, "Internal Server Error")
		return
	}

	//set cookie using generated JWT and send toast message to inform user of correct login

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   signedString,
		Path:    "/",
		Expires: expirationTime,
	})
	w.Header().Add("HX-Redirect", "/notes")
	w.WriteHeader(http.StatusOK)
}

func (app *App) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Redirect", "/")
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true})
	w.WriteHeader(200)
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
