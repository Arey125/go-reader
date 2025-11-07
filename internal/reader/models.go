package reader

import (
	"database/sql"
	"reader/internal/dictionary"
)

type TextModel interface {
	Add(text Text) error
	AllWithoutContent() ([]Text, error)
	Get(id int) (*Text, error)
	GetPage(textId int, page int) (*TextPage, error)
}

type WordModel interface {
	BeginTx() (*sql.Tx, WordModel, error)
	AddList(words []Word) error
	SaveDefinitions(word Word, definitions []dictionary.Definition) error
	GetDefinitions(word Word) ([]dictionary.Definition, error)
	AddUserWordList(words []Word, userId int) error
	GetUserWords(userId int) ([]UserWord, error)
}
