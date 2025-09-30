package reader

import (
	"net/http"
	"reader/internal/dictionary"
	"reader/internal/nlp"
	"reader/internal/server"
	"reader/internal/users"
	"strconv"
	"strings"
)

type Service struct {
	model            *Model
	nlpClient        *nlp.Client
	dictionaryClient *dictionary.Client
	wordFreq WordFreq
}

func NewService(model *Model, nlpClient *nlp.Client, dictionaryClient *dictionary.Client) Service {
	return Service{
		model:            model,
		nlpClient:        nlpClient,
		dictionaryClient: dictionaryClient,
		wordFreq: NewWordFreq(),
	}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", s.homePage)
	mux.HandleFunc("GET /texts/add", s.addPage)
	mux.HandleFunc("POST /texts/add", s.addPost)
	mux.HandleFunc("GET /texts/{id}", s.readPage)

	mux.HandleFunc("GET /word", s.wordGet)
}

func (s *Service) homePage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	texts, err := s.model.All()
	if err != nil {
		server.ServerError(w, err)
	}
	homePageTempl(user, texts).Render(r.Context(), w)
}

func (s *Service) readPage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	textIdStr := r.PathValue("id")
	pageIndStr := r.FormValue("page")
	if pageIndStr == "" {
		pageIndStr = "1"
	}

	textId, err := strconv.Atoi(textIdStr)
	if err != nil {
		server.HttpError(w, http.StatusBadRequest)
		return
	}
	pageInd, err := strconv.Atoi(pageIndStr)
	if err != nil {
		server.HttpError(w, http.StatusBadRequest)
		return
	}

	pagePtr, err := s.model.GetPage(textId, pageInd)
	if err != nil {
		server.ServerError(w, err)
	}
	if pagePtr == nil {
		server.HttpError(w, http.StatusNotFound)
		return
	}
	page := *pagePtr
	page.Content = strings.TrimSpace(page.Content)

	segments, err := s.splitIntoSegments(page.Content)
	if err != nil {
		server.ServerError(w, err)
	}

	readPageTempl(user, page, segments).Render(r.Context(), w)
}

func (s *Service) wordGet(w http.ResponseWriter, r *http.Request) {
	segment := Segment{}
	segment.Text = r.FormValue("text")
	segment.Info = &WordInfo{
		Lemma: r.FormValue("lemma"),
		Pos:   r.FormValue("pos"),
	}

	definitions, err := s.getDefinitions(*segment.Info)
	if err != nil {
		server.ServerError(w, err)
	}

	freq := s.wordFreq.Get(segment.Info.Lemma)
	wordTempl(segment, definitions, freq).Render(r.Context(), w)
}
