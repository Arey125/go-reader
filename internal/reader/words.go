package reader

func (s *Service) SaveWordsFromSegments(segments []Segment) error {
	words := getWordsFromSegments(segments)

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

func (s *Service) SaveWordsFromSegmentsAsKnown(segments []Segment, userId int) error {
	words := getWordsFromSegments(segments)

	tx, wordModelWithTx, err := s.wordModel.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = wordModelWithTx.AddUserWordList(words, userId)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func getWordsFromSegments(segments []Segment) []Word {
	words := make([]Word, 0)
	for _, segment := range segments {
		if segment.Info == nil {
			continue
		}

		words = append(words, Word{
			Word: segment.Info.Lemma,
			Pos:  segment.Info.Pos,
		})
	}
	return words
}
