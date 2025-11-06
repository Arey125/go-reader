package reader

func (s *Service) SaveWordsFromPage(page TextPage) error {
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

	tx, wordModelWithTx, err := s.wordModel.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = wordModelWithTx.AddList(words)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}
