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

type TextPage struct {
	TextId int
	TextTitle string
	Page int
	Total int
	Content string
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

func (m *Model) Get(id int) (*Text, error) {
	t := Text{}

	err := sq.Select("id", "title", "content", "user_id").
		From("texts").
		Where(sq.Eq{"id": id}).
		RunWith(m.db).
		QueryRow().
		Scan(&t.Id, &t.Title, &t.Content, &t.UserId)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (m *Model) GetPage(textId int, page int) (*TextPage, error) {
	t := Text{}

	err := sq.Select("id", "title", "content", "user_id").
		From("texts").
		Where(sq.Eq{"id": textId}).
		RunWith(m.db).
		QueryRow().
		Scan(&t.Id, &t.Title, &t.Content, &t.UserId)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	pages := splitIntoPages(t.Content, 1000)
	if len(pages) < page {
		return nil, nil
	}
	p := TextPage{
		TextId: t.Id,
		TextTitle: t.Title,
		Content: pages[page - 1],
		Page: page,
		Total: len(pages),
	}

	return &p, nil
}
