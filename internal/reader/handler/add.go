package handler

import (
	"io"
	"net/http"
	"reader/internal/reader"
	"reader/internal/server"
	"reader/internal/users"
	"time"
)

func (s *Handler) addPage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	addPageTempl(user).Render(r.Context(), w)
}

func (s *Handler) addPost(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	title := r.FormValue("title")
	uploadType := r.FormValue("type")

	if user == nil || len(title) < 3 {
		server.HttpError(w, http.StatusBadRequest)
		return
	}

	content := ""
	if uploadType == "text" {
		content = r.FormValue("content")
	}
	if uploadType == "txt" {
		file, _, err := r.FormFile("file")
		if err != nil {
			server.ServerError(w, err)
			return
		}
		contentBytes, err := io.ReadAll(file)
		if err != nil {
			server.ServerError(w, err)
			return
		}
		content = string(contentBytes)
	}

	if len(content) < 3 {
		server.HttpError(w, http.StatusBadRequest)
		return
	}

	err := s.service.AddText(reader.Text{
		Title:   title,
		Content: content,
		UserId:  user.User.Id,
		CreatedAt: time.Now(),
	})

	if err != nil {
		server.ServerError(w, err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
