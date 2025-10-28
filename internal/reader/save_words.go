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

	tx, err := s.wordModel.WithTx()
	defer tx.Rollback()

	if err != nil {
		return err
	}

	err = tx.AddList(words)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
