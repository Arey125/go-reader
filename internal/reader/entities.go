package reader

import (
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

type UserWord struct {
	Word Word
	Status string
}
