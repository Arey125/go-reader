package reader

import (
	"database/sql"
	"encoding/json"
	"reader/internal/dictionary"

	sq "github.com/Masterminds/squirrel"
)

type Word struct {
	Id          int
	Word        string
	Pos         string
	Definitions *string
}

type WordModel struct {
	db *sql.DB
}

func NewWordModel(db *sql.DB) WordModel {
	return WordModel{db}
}

func (m *WordModel) AddList(words []Word) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	for _, word := range words {
		_, err := sq.Insert("").
			Options("OR IGNORE").
			Into("words").
			Columns("word", "pos").
			Values(word.Word, word.Pos).
			RunWith(tx).
			Exec()

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (m *WordModel) SaveDefinitions(word Word, definitions []dictionary.Definition) error {
	definitionsJson, err := json.Marshal(definitions)
	if err != nil {
		return err
	}

	definitionsStr := string(definitionsJson)

	_, err = sq.Update("words").
		Set("definitions", definitionsStr).
		Where(sq.Eq{"word": word.Word, "pos": word.Pos}).
		RunWith(m.db).
		Exec()

	return err
}
