package reader

import "reader/internal/dictionary"

type TextModel interface {
	Add(text Text) error
	AllWithoutContent() ([]Text, error)
	Get(id int) (*Text, error)
	GetPage(textId int, page int) (*TextPage, error)
}

type WordModel interface {
	AddList(words []Word) error
	SaveDefinitions(word Word, definitions []dictionary.Definition) error
	GetDefinitions(word Word) ([]dictionary.Definition, error)
}
