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
	db sq.BaseRunner
}

func NewWordModel(db *sql.DB) WordModel {
	return WordModel{db}
}

func (m *WordModel) WithTx() (*WordModel, error) {
	tx, err := m.db.(*sql.DB).Begin()
	if err != nil {
		return nil, err
	}
	return &WordModel{tx}, nil
}

func (m *WordModel) Commit() error {
	return m.db.(*sql.Tx).Commit()
}

func (m *WordModel) Rollback() error {
	return m.db.(*sql.Tx).Rollback()
}

func (m *WordModel) AddList(words []Word) error {
	for _, word := range words {
		_, err := sq.Insert("").
			Options("OR IGNORE").
			Into("words").
			Columns("word", "pos").
			Values(word.Word, word.Pos).
			RunWith(m.db).
			Exec()

		if err != nil {
			return err
		}
	}
	return nil
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

func (m *WordModel) GetDefinitions(word Word) ([]dictionary.Definition, error) {
	var definitionsJsonStr *string

	err := sq.Select("definitions").
		From("words").
		Where(sq.Eq{"word": word.Word, "pos": word.Pos}).
		RunWith(m.db).
		QueryRow().
		Scan(&definitionsJsonStr)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if definitionsJsonStr == nil {
		return nil, nil
	}
	
	definitions := []dictionary.Definition{}
	err = json.Unmarshal([]byte(*definitionsJsonStr), &definitions)
	if err != nil {
		return nil, err
	}

	return definitions, nil
}
