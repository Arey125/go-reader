package handler

import (
	"net/http"
	"reader/internal/reader"
	"reader/internal/server"
	"reader/internal/users"
	"strconv"

	"reader/internal/reader/ui"

	"github.com/go-playground/form/v4"
)

type readPageQuery struct {
	Text int
	Page int
	From *int
}

func (h *Handler) readPage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	query, err := getReadPageQuery(w, r)
	if err != nil {
		return
	}

	page, segments, err := h.getPageAndSegments(w, query.Text, query.Page)
	if err != nil {
		return
	}
	err = h.service.SaveWordsFromSegments(segments)
	if err != nil {
		server.ServerError(w, err)
		return
	}

	if query.From != nil {
		_, fromSegments, err := h.getPageAndSegments(w, query.Text, *query.From)
		if err != nil {
			return
		}
		err = h.service.SaveWordsFromSegmentsAsKnown(fromSegments, user.User.Id)
		if err != nil {
			server.ServerError(w, err)
			return
		}
	}

	ui.ReadPage(user, page, segments).Render(r.Context(), w)
}

type readPageForm struct {
	Page *int `form:"page,omitempty"`
	From *int `form:"from,omitempty"`
}

func getReadPageQuery(w http.ResponseWriter, r *http.Request) (readPageQuery, error) {
	textId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		server.HttpError(w, http.StatusBadRequest)
		return readPageQuery{}, err
	}

	formDecoder := form.NewDecoder()
	form := readPageForm{}
	err = r.ParseForm()
	if err != nil {
		server.HttpError(w, http.StatusBadRequest)
		return readPageQuery{}, err
	}

	err = formDecoder.Decode(&form, r.Form)
	if err != nil {
		server.HttpError(w, http.StatusBadRequest)
		return readPageQuery{}, err
	}

	pageInd := 1
	if form.Page != nil {
		pageInd = *form.Page
	}

	return readPageQuery{
		Text: textId,
		Page: pageInd,
		From: form.From,
	}, nil
}

func (h *Handler) getPageAndSegments(
	w http.ResponseWriter,
	textId int,
	pageInd int,
) (reader.TextPage, []reader.Segment, error) {
	pagePtr, segments, err := h.service.GetPageAndSegments(textId, pageInd)
	if err != nil {
		server.ServerError(w, err)
		return reader.TextPage{}, nil, err
	}
	if pagePtr == nil {
		server.HttpError(w, http.StatusNotFound)
		return reader.TextPage{}, nil, err
	}
	page := *pagePtr

	return page, segments, err
}
