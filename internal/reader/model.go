package reader

import (
	"database/sql"
	"reader/internal/db"

	sq "github.com/Masterminds/squirrel"
)

type Text struct {
	Id      int
	Title   string
	Content string
	UserId  int
}

type Model struct {
	db *sql.DB
}

func NewModel(db *sql.DB) Model {
	return Model{db}
}

func (m *Model) Add(text Text) error {
	_, err := sq.Insert("texts").
		Columns("title", "content", "user_id").
		Values(text.Title, text.Content, text.UserId).
		RunWith(m.db).
		Exec()

	return err
}

func (m *Model) All() ([]Text, error) {
	rows, err := sq.Select("id", "title", "content", "user_id").
		From("texts").
		RunWith(m.db).
		Query()

	if err != nil {
		return nil, err
	}

	return db.Collect(rows, func(r *sql.Rows, t *Text) error {
		return rows.Scan(&t.Id, &t.Title, &t.Content, &t.UserId)
	})
}
