package word

import (
	"context"
	"reader/internal/db/queries"

	"reader/internal/reader"
)

func (m *WordModel) AddUserWordList(words []reader.Word, userId int) error {
	for _, word := range words {
		err := m.q.AddUserWord(context.Background(), queries.AddUserWordParams{
			UserID: int64(userId),
			Word: word.Word,
			Pos: word.Pos,
			Status: "known",
		})

		if err != nil {
			return err
		}
	}
	return nil
}
