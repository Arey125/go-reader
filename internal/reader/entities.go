package reader

import (
	"fmt"
	"time"
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

type Word struct {
	Id          int
	Word        string
	Pos         string
	Definitions *string
}

type WordStatus string

const (
	WordStatusKnown    WordStatus = "known"
	WordStatusLearning WordStatus = "learning"
	WordStatusLearned  WordStatus = "learned"
)

type UserWord struct {
	Word   Word
	Status WordStatus
}

func ToWordStatus(status string) (WordStatus, error) {
	switch status {
	case string(WordStatusKnown):
		return WordStatusKnown, nil
	case string(WordStatusLearning):
		return WordStatusLearning, nil
	case string(WordStatusLearned):
		return WordStatusLearned, nil
	default:
		return "", fmt.Errorf("unexpected reader.WordStatus: %#v", status)
	}
}
