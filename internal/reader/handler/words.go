package handler

import (
	"net/http"
	"reader/internal/reader"
	"reader/internal/server"
	"reader/internal/users"
	"slices"

	"reader/internal/reader/ui"
)

func (h *Handler) wordsPage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	words, err := h.service.GetUserWordsWithFreq(user.User.Id)
	if err != nil {
		server.ServerError(w, err)
	}
	slices.SortFunc(words, func(a reader.UserWordWithFreq, b reader.UserWordWithFreq) int {
		return b.Freq - a.Freq
	})
	ui.WordsPage(user, words).Render(r.Context(), w)
}
