package reader

import (
	"reader/internal/dictionary"
)

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

	return s.dictionaryClient.GetDefinitions(wordInfo.Lemma, wordInfo.Pos)
}
