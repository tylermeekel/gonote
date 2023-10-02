package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct{
	ID int
	Username string
	Password []byte
}

func (app App) userRouter() *chi.Mux{
	router := chi.NewRouter()
	router.Post("/", app.handleCreateUser)

	return router
}

func (app App) authRouter() *chi.Mux{
	router := chi.NewRouter()
	router.Post("/login", app.handleLoginUser)

	return router
}

func (app App) getUserByUsername(username string) (User, error){
	var user User
	row, err := app.db.Query("SELECT * FROM users WHERE username=$1", username)
	if err != nil{
		return user, err
	}
	if row.Next(){
		err = row.Scan(&user.ID, &user.Username, &user.Password)
	}
	return user, err
}

func (app App) createUser(username string, password []byte) (User, error){
	var user User
	row, err := app.db.Query("INSERT INTO users(username, password) VALUES($1, $2) RETURNING *", username, password)
	if err != nil{
		return user, err
	}
	if row.Next(){
		row.Scan(&user.ID, &user.Username, &user.Password)
	}

	return user, nil
}

func (app App) handleCreateUser(w http.ResponseWriter, r *http.Request){
	username := r.FormValue("username")
	password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 10)
	if err != nil{
		fmt.Println(err.Error())
	}

	user, err := app.createUser(username, password)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}

	//TODO: change later
	json.NewEncoder(w).Encode(user)
}

func (app App) handleLoginUser(w http.ResponseWriter, r *http.Request){
	username := r.FormValue("username")
	password := r.FormValue("password")

	queriedUser, err := app.getUserByUsername(username)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(queriedUser.Password), []byte(password))
	if err != nil || queriedUser.Username == ""{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request: either username or password were incorrect"))
		//TODO: store JWT in cookie
	} else {
		w.Write([]byte("correct password"))
	}
}