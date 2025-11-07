package dictionary

import (
	"slices"
)

func getDefinitionsFromEntry(entry DictionaryEntry, wordPos string) []Definition {
	meanings := entry.Meanings

	toPos := map[string]string{
		"adjective":    "ADJ",
		"adverb":       "ADV",
		"interjection": "INTJ",
		"noun":         "NOUN",
		"verb":         "VERB",
	}
	foundInd := slices.IndexFunc(meanings, func(meaning Meaning) bool {
		pos, ok := toPos[meaning.Pos]
		if !ok {
			return false
		}
		return pos == wordPos
	})

	if foundInd == -1 {
		return []Definition{{}}
	}

	return meanings[foundInd].Definitions
}

func (c *Client) GetDefinitions(word string, pos string) ([]Definition, error) {
	entries, err := c.GetEntries(word)

	if err != nil {
		return nil, err
	}

	definitions := make([]Definition, 0)
	for _, entry := range entries {
		definitions = append(definitions, getDefinitionsFromEntry(entry, pos)...)
	}

	return definitions, nil
}
