package reader

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Word struct {
	Id         int
	Word       string
	Pos        string
	Definition *string
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
