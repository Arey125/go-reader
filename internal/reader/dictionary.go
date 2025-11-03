package reader

import (
	"reader/internal/dictionary"
	"slices"
)

func getDefinitionsFromEntry(entry dictionary.DictionaryEntry, wordPos string) []dictionary.Definition {
	meanings := entry.Meanings

	toPos := map[string]string{
		"adjective":    "ADJ",
		"adverb":       "ADV",
		"interjection": "INTJ",
		"noun":         "NOUN",
		"verb":         "VERB",
	}
	foundInd := slices.IndexFunc(meanings, func(meaning dictionary.Meaning) bool {
		pos, ok := toPos[meaning.Pos]
		if !ok {
			return false
		}
		return pos == wordPos
	})

	if foundInd == -1 {
		return []dictionary.Definition{{}}
	}

	return meanings[foundInd].Definitions
}

func (s *Service) getDefinitions(wordInfo WordInfo) ([]dictionary.Definition, error) {
	word := Word{
		Word: wordInfo.Lemma,
		Pos:  wordInfo.Pos,
	}
	definitions, err := s.wordModel.GetDefinitions(word)
	if err != nil {
		return nil, err
	}
	if definitions != nil {
		return definitions, nil
	}

	entries, err := s.dictionaryClient.GetEntries(wordInfo.Lemma)

	if err != nil {
		return nil, err
	}

	definitions = make([]dictionary.Definition, 0)
	for _, entry := range entries {
		definitions = append(definitions, getDefinitionsFromEntry(entry, wordInfo.Pos)...)
	}

	return definitions, nil
}
