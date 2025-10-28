package reader

import (
	"strings"
	"errors"
)

var PageNotFoundErr error = errors.New("page not found")

func (s *Service) getPage(textId int, pageInd int) (*TextPage, error) {
	pagePtr, err := s.textModel.GetPage(textId, pageInd)
	if err != nil {
		return nil, err
	}
	if pagePtr == nil {
		return nil, nil
	}

	pagePtr.Content = strings.TrimSpace(pagePtr.Content)
	return pagePtr, nil
}
