package main

import (
	"fmt"
)

type Note struct {
	title   string
	content string
}

func (app App) getAllNotes() []Note {
	var notes []Note
	rows, err := app.db.Query("SELECT * FROM notes;")
	if err != nil {
		fmt.Println("error querying notes")
	}
	for rows.Next() {
		var note Note
		rows.Scan(&note.title, &note.content)
		notes = append(notes, note)
	}

	return notes
}
