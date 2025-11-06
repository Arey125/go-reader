package reader

import (
	"context"
	"database/sql"
	"reader/internal/db/queries"
	"reader/internal/reader"

	sq "github.com/Masterminds/squirrel"

)

type TextModel struct {
	db sq.BaseRunner
	q  *queries.Queries
}

func New(db *sql.DB, q *queries.Queries) TextModel {
	return TextModel{db, q}
}

func (m *TextModel) Add(text reader.Text) error {
	return m.q.AddText(context.Background(), queries.AddTextParams{
		Title: text.Title,
		Content: text.Content,
		UserID: int64(text.UserId),
		CreatedAt: text.CreatedAt,
	})
}

func (m *TextModel) AllWithoutContent() ([]reader.Text, error) {
	rows, err := m.q.AllTextsWithoutContent(context.Background())
	if err != nil {
		return nil, err
	}
	texts := make([]reader.Text, len(rows))

	for i, row := range rows {
		texts[i] = reader.Text{
			Id:        int(row.ID),
			Title:     row.Title,
			UserId:    int(row.UserID),
			CreatedAt: row.CreatedAt,
		}
	}
	return texts, nil
}

func (m *TextModel) Get(id int) (*reader.Text, error) {
	row, err := m.q.GetText(context.Background(), int64(id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &reader.Text{
		Id:        int(row.ID),
		Title:     row.Title,
		Content:   row.Content,
		UserId:    int(row.UserID),
		CreatedAt: row.CreatedAt,
	}, nil
}

func (m *TextModel) GetPage(textId int, page int) (*reader.TextPage, error) {
	t, err := m.Get(textId)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	pages := reader.SplitIntoPages(t.Content, 1000)
	if len(pages) < page {
		return nil, nil
	}
	p := reader.TextPage{
		TextId:    t.Id,
		TextTitle: t.Title,
		Content:   pages[page-1],
		Page:      page,
		Total:     len(pages),
	}

	return &p, nil
}
