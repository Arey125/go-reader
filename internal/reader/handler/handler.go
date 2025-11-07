package handler

import (
	"net/http"
	"reader/internal/server"
	"reader/internal/users"

	"reader/internal/reader"
	"reader/internal/reader/ui"
)

type Handler struct {
	service *reader.Service
}

func New(service *reader.Service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", h.homePage)
	mux.HandleFunc("GET /texts/add", h.addPage)
	mux.HandleFunc("POST /texts/add", h.addPost)
	mux.HandleFunc("GET /texts/{id}", h.readPage)

	mux.HandleFunc("GET /word", h.wordGet)
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	texts, err := h.service.GetAllTexts()
	if err != nil {
		server.ServerError(w, err)
	}
	ui.HomePage(user, texts).Render(r.Context(), w)
}


func (h *Handler) wordGet(w http.ResponseWriter, r *http.Request) {
	segment := reader.Segment{}
	segment.Text = r.FormValue("text")
	segment.Info = &reader.WordInfo{
		Lemma: r.FormValue("lemma"),
		Pos:   r.FormValue("pos"),
	}

	definitions, freq, err := h.service.GetWordDefinitionsAndFreq(segment)
	if err != nil {
		server.ServerError(w, err)
	}

	ui.Word(segment, definitions, freq).Render(r.Context(), w)
}
