package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"math"
	"net/http"
)

type Sharelink struct {
	ID          string
	Title       string
	Content     string
	ContentHTML template.HTML
}

func (app *App) sharelinkRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", app.handleCreateSharelink)
	r.Get("/{id}", app.handleGetSharelink)

	return r
}

func (app *App) createSharelink(title string, content string) string {

	buff := make([]byte, int(math.Ceil(float64(32)/2)))
	_, err := rand.Read(buff)
	if err != nil {
		app.log.Println("Error creating random string for sharelink: ", err.Error())
	}
	str := hex.EncodeToString(buff)

	row := app.db.QueryRow("INSERT INTO share_links(id, title, content) VALUES($1, $2, $3) RETURNING *", str, title, content)
	err = row.Scan()
	if err != nil {
		app.log.Println("Error creating sharelink", err.Error())
	}

	return str[:32]
}

func (app *App) getSharelinkContent(id string) Sharelink {
	row := app.db.QueryRow("SELECT * FROM share_links WHERE id=$1", id)
	var sharelink Sharelink
	err := row.Scan(&sharelink.ID, &sharelink.Title, &sharelink.Content)
	if err != nil {
		app.log.Println("Error getting sharelink", err.Error())
	}

	return sharelink
}

func (app *App) handleGetSharelink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	note := app.getSharelinkContent(id)
	note.ContentHTML = template.HTML(mdToHTML(note.Content))

	err := app.templates.ExecuteTemplate(w, "sharelink", note)
	if err != nil {
		app.log.Println("Error executing sharelink template: ", err.Error())
	}
}

func (app *App) handleCreateSharelink(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")

	app.log.Println(title, content)

	sharelink := app.createSharelink(title, content)
	redirect := fmt.Sprintf("/sharelink/%s", sharelink)
	w.Header().Set("HX-Redirect", redirect)
	w.WriteHeader(200)
}
