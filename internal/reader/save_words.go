package reader

func (s *Service) saveWordsFromPage(page TextPage) error {
	segments, err := s.nlpClient.GetWords(page.Content)
	if err != nil {
		return err
	}

	words := make([]Word, 0, len(segments))
	for _, segment := range segments {
		words = append(words, Word{
			Word: segment.Lemma,
			Pos:  segment.Pos,
		})
	}

	err = s.wordModel.AddList(words)
	if err != nil {
		return err
	}

	return nil
}
