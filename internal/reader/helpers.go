package reader

import (
	"net/http"
	"reader/internal/server"
	"strconv"
	"strings"
	"errors"
)

var PageNotFoundErr error = errors.New("page not found")

func (s *Service) getPage(w http.ResponseWriter, textIdStr string, pageIndStr string) (*TextPage, error) {
	if pageIndStr == "" {
		pageIndStr = "1"
	}

	textId, err := strconv.Atoi(textIdStr)
	if err != nil {
		server.HttpError(w, http.StatusBadRequest)
		return nil, err
	}

	pageInd, err := strconv.Atoi(pageIndStr)
	if err != nil {
		server.HttpError(w, http.StatusBadRequest)
		return nil, err
	}

	pagePtr, err := s.textModel.GetPage(textId, pageInd)
	if err != nil {
		server.ServerError(w, err)
	}
	if pagePtr == nil {
		server.HttpError(w, http.StatusNotFound)
		return nil, PageNotFoundErr
	}

	pagePtr.Content = strings.TrimSpace(pagePtr.Content)
	return pagePtr, nil
}
