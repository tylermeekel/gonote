package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Password []byte
}

func (app *App) userRouter() *chi.Mux {
	router := chi.NewRouter()

	return router
}

func (app *App) authRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/register", app.handleRegisterUser)
	router.Post("/login", app.handleLoginUser)

	return router
}

func (app *App) getUserByUsername(username string) (User, error) {
	var user User
	row, err := app.db.Query("SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		return user, err
	}
	if row.Next() {
		err = row.Scan(&user.ID, &user.Username, &user.Password)
	}
	return user, err
}

func (app *App) createUser(username string, password []byte) (User, error) {
	var user User
	row, err := app.db.Query("INSERT INTO users(username, password) VALUES($1, $2) RETURNING *", username, password)
	if err != nil {
		return user, err
	}
	if row.Next() {
		row.Scan(&user.ID, &user.Username, &user.Password)
	}

	return user, nil
}

func (app *App) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 10)
	if err != nil {
		fmt.Println(err.Error())
	}

	user, err := app.createUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErrorToast(w, "Error: Internal Server Error")
	} else {
		fmt.Printf("Created user \"%s\"\n", user.Username)
		sendToast(w, "User successfully created")
	}

}

func (app *App) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	isSuccessful := true
	username := r.FormValue("username")
	password := r.FormValue("password")

	queriedUser, err := app.getUserByUsername(username)
	if err != nil {
		sendErrorToast(w, "Error: Internal Server Error")
		isSuccessful = false
	}

	//Check hash and password and return an error if they do not match
	err = bcrypt.CompareHashAndPassword([]byte(queriedUser.Password), []byte(password))
	if err != nil || queriedUser.Username == "" {
		sendErrorToast(w, "Incorrect username or password")
		return
	}

	//set expiration time and create JWT using username and expiration time and return error if there is
	//an error generating the jwt
	expirationTime := time.Now().Add(10 * time.Minute)
	signedString, err := signJWT(username, expirationTime)
	if err != nil {
		app.log.Println(err.Error())
		sendErrorToast(w, "Internal Server Error")
		return
	}

	//set cookie using generated JWT and send toast message to inform user of correct login
	if isSuccessful {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    signedString,
			Expires:  expirationTime,
		})
		sendToast(w, "Logged in successfully")
	}
}
