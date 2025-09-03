package reader

type WordInfo struct {
	Lemma string
	Pos   string
}

type Segment struct {
	Text string
	Info *WordInfo
}

func (s *Service) splitIntoSegments(text string) ([]Segment, error) {
	words, err := s.nlpClient.GetWords(text)
	if err != nil {
		return nil, err
	}

	textRunes := []rune(text)
	textInd := 0
	wordInd := 0
	res := make([]Segment, 0)
	for textInd != len(textRunes) {
		if wordInd == len(words) {
			res = append(res, Segment{string(textRunes[textInd:]), nil})
			break
		}
		nextWord := words[wordInd]
		if textInd < nextWord.Start {
			res = append(res, Segment{string(textRunes[textInd:nextWord.Start]), nil})
		}
		res = append(res, Segment{nextWord.Text, &WordInfo{nextWord.Lemma, nextWord.Pos}})

		textInd = nextWord.End
		wordInd++
	}

	return res, nil
}
