package word

import (
	"context"
	"reader/internal/db/queries"

	"reader/internal/reader"
)

func (m *WordModel) AddUserWordListAsKnown(words []reader.Word, userId int) error {
	for _, word := range words {
		err := m.q.AddOrIgnoreUserWord(context.Background(), queries.AddOrIgnoreUserWordParams{
			UserID: int64(userId),
			Word:   word.Word,
			Pos:    word.Pos,
			Status: string(reader.WordStatusKnown),
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (m *WordModel) AddOrReplaceUserWordAsLearning(word reader.Word, userId int) error {
	return m.q.AddOrReplaceUserWord(context.Background(), queries.AddOrReplaceUserWordParams{
		UserID: int64(userId),
		Word:   word.Word,
		Pos:    word.Pos,
		Status: string(reader.WordStatusLearning),
	})
}

func (m *WordModel) GetUserWords(userId int) ([]reader.UserWord, error) {
	dbWords, err := m.q.GetUserWords(context.Background(), int64(userId))
	if err != nil {
		return nil, err
	}
	words := make([]reader.UserWord, len(dbWords))
	for i, dbWord := range dbWords {
		status, err := reader.ToWordStatus(dbWord.Status)
		if err != nil {
			return nil, err
		}

		words[i] = reader.UserWord{
			Word: reader.Word{
				Id:   int(dbWord.ID),
				Word: dbWord.Word,
				Pos:  dbWord.Pos,
			},
			Status: status,
		}
	}
	return words, nil
}
