package main

import (
	"fmt"
)

type Note struct {
	Title   string
	Content string
}

func (app App) getAllNotes() []Note {
	var notes []Note
	rows, err := app.db.Query("SELECT * FROM notes;")
	if err != nil {
		fmt.Println("error querying notes")
	}
	for rows.Next() {
		var note Note
		rows.Scan(&note.Title, &note.Content)
		notes = append(notes, note)
	}

	return notes
}
