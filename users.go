package main

type User struct {
	ID       int
	Username string
	Password []byte
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
func (app *App) createUser(username ValidUsername, passwordHash []byte) (User, error) {
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
