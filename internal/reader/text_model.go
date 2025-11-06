package reader

import (
	"context"
	"database/sql"
	"reader/internal/db/queries"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Text struct {
	Id        int
	Title     string
	Content   string
	UserId    int
	CreatedAt time.Time
}

type TextPage struct {
	TextId    int
	TextTitle string
	Page      int
	Total     int
	Content   string
}

type TextModel struct {
	db sq.BaseRunner
	q  *queries.Queries
}

func NewTextModel(db *sql.DB, q *queries.Queries) TextModel {
	return TextModel{db, q}
}

func (m *TextModel) Add(text Text) error {
	return m.q.AddText(context.Background(), queries.AddTextParams{
		Title: text.Title,
		Content: text.Content,
		UserID: int64(text.UserId),
		CreatedAt: text.CreatedAt,
	})
}

func (m *TextModel) AllWithoutContent() ([]Text, error) {
	rows, err := m.q.AllTextsWithoutContent(context.Background())
	if err != nil {
		return nil, err
	}
	texts := make([]Text, len(rows))

	for i, row := range rows {
		texts[i] = Text{
			Id:        int(row.ID),
			Title:     row.Title,
			UserId:    int(row.UserID),
			CreatedAt: row.CreatedAt,
		}
	}
	return texts, nil
}

func (m *TextModel) Get(id int) (*Text, error) {
	row, err := m.q.GetText(context.Background(), int64(id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &Text{
		Id:        int(row.ID),
		Title:     row.Title,
		Content:   row.Content,
		UserId:    int(row.UserID),
		CreatedAt: row.CreatedAt,
	}, nil
}

func (m *TextModel) GetPage(textId int, page int) (*TextPage, error) {
	t, err := m.Get(textId)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	pages := splitIntoPages(t.Content, 1000)
	if len(pages) < page {
		return nil, nil
	}
	p := TextPage{
		TextId:    t.Id,
		TextTitle: t.Title,
		Content:   pages[page-1],
		Page:      page,
		Total:     len(pages),
	}

	return &p, nil
}
