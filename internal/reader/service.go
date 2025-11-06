package reader

import (
	"errors"
	"reader/internal/dictionary"
	"reader/internal/nlp"
)

type Service struct {
	textModel        TextModel
	wordModel        WordModel
	nlpClient        *nlp.Client
	dictionaryClient *dictionary.Client
	wordFreq         WordFreq
}

func NewService(
	textModel TextModel,
	wordModel WordModel,
	nlpClient *nlp.Client,
	dictionaryClient *dictionary.Client,
) Service {
	return Service{
		textModel:        textModel,
		wordModel:        wordModel,
		nlpClient:        nlpClient,
		dictionaryClient: dictionaryClient,
		wordFreq:         NewWordFreq(),
	}
}

var ErrNotFound = errors.New("Not found")

func (s *Service) GetAllTexts() ([]Text, error) {
	return s.textModel.AllWithoutContent()
}

func (s *Service) GetPageAndSegments(textId int, pageInd int) (*TextPage, []Segment, error) {
	pagePtr, err := s.getPage(textId, pageInd)
	if err != nil {
		return nil, nil, err
	}
	if pagePtr == nil {
		return nil, nil, nil
	}

	page := *pagePtr
	segments, err := s.splitIntoSegments(page.Content)
	if err != nil {
		return nil, nil, err
	}
	if err = s.SaveWordsFromPage(page); err != nil {
		return nil, nil, err
	}

	return pagePtr, segments, nil
}

func (s *Service) GetWordDefinitionsAndFreq(segment Segment) ([]dictionary.Definition, *WordFreqRecord, error) {
	definitions, err := s.getDefinitions(*segment.Info)
	if err != nil {
		return nil, nil, err
	}

	word := Word{
		Word: segment.Info.Lemma,
		Pos:  segment.Info.Pos,
	}
	s.wordModel.AddList([]Word{word})
	s.wordModel.SaveDefinitions(Word{
		Word: segment.Info.Lemma,
		Pos:  segment.Info.Pos,
	}, definitions)

	freq := s.wordFreq.Get(segment.Info.Lemma)
	return definitions, freq, nil
}

func (s *Service) AddText(text Text) error {
	return s.textModel.Add(text)
}
