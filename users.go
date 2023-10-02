package main

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct{
	ID int
	Username string
	Password []byte
}

func (app App) createUser(username string, password []byte) error{
	_, err := app.db.Query("INSERT INTO users(username, password) VALUES($1, $2)", username, password)
	if err != nil{
		return err
	}

	return nil
}

func (app App) handleCreateUser(w http.ResponseWriter, r *http.Request){
	username := r.FormValue("username")
	password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 10)
	if err != nil{
		fmt.Println(err.Error())
	}

	err = app.createUser(username, password)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}
}