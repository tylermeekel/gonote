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

// func (app *App) userRouter() *chi.Mux {
// 	router := chi.NewRouter()

// 	return router
// }

func (app *App) authRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/register", app.handleRegisterUser)
	router.Post("/login", app.handleLoginUser)

	return router
}

// getUserByUsername queries the database for a username and returns a User struct with the
// information gathered
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

// createUser takes a username and password hash as input and inserts them into the database
func (app *App) createUser(username string, passwordHash []byte) (User, error) {
	var user User
	row, err := app.db.Query("INSERT INTO users(username, password) VALUES($1, $2) RETURNING *", username, passwordHash)
	if err != nil {
		return user, err
	}
	if row.Next() {
		row.Scan(&user.ID, &user.Username, &user.Password)
	}

	return user, nil
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
		sendErrorToast(w, "Error: Internal Server Error")
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
	if err != nil {
		sendErrorToast(w, "Error: Internal Server Error")
		return
	}

	//Check hash and password and return an error if they do not match
	err = bcrypt.CompareHashAndPassword([]byte(queriedUser.Password), []byte(password))
	if err != nil || queriedUser.Username == "" {
		sendErrorToast(w, "Incorrect username or password")
		return
	}

	//set expiration time and create JWT using user ID and expiration time and return error if there is
	//an error generating the jwt
	expirationTime := time.Now().Add(10 * time.Minute)
	signedString, err := signJWT(queriedUser.ID, expirationTime)
	if err != nil {
		app.log.Println(err.Error())
		sendErrorToast(w, "Internal Server Error")
		return
	}

	//set cookie using generated JWT and send toast message to inform user of correct login

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   signedString,
		Expires: expirationTime,
	})
	sendToast(w, "Logged in successfully")

}

