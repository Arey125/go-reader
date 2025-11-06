package word

import (
	"context"
	"database/sql"
	"encoding/json"
	"reader/internal/db/queries"
	"reader/internal/dictionary"

	"reader/internal/reader"
)

type WordModel struct {
	db *sql.DB
	q  *queries.Queries
}

func New(db *sql.DB, q *queries.Queries) WordModel {
	return WordModel{db, q}
}

func (m *WordModel) AddList(words []reader.Word) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := m.q.WithTx(tx)
	for _, word := range words {
		err := qtx.AddWord(context.Background(), queries.AddWordParams{
			Word: word.Word,
			Pos: word.Pos,
		})

		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (m *WordModel) SaveDefinitions(word reader.Word, definitions []dictionary.Definition) error {
	definitionsJson, err := json.Marshal(definitions)
	if err != nil {
		return err
	}

	definitionsStr := string(definitionsJson)

	return m.q.UpdateWordDefinition(context.Background(), queries.UpdateWordDefinitionParams{
		Definitions: &definitionsStr,
		Pos: word.Pos,
		Word: word.Word,
	})
}

func (m *WordModel) GetDefinitions(word reader.Word) ([]dictionary.Definition, error) {
	definitionsJsonStr, err := m.q.GetWordDefinition(context.Background(), queries.GetWordDefinitionParams{
		Word: word.Word,
		Pos: word.Pos,
	})

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
