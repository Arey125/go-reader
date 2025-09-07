package reader

import (
	"net/http"
	"reader/internal/dictionary"
	"reader/internal/nlp"
	"reader/internal/server"
	"reader/internal/users"
	"strconv"
)

type Service struct {
	model            *Model
	nlpClient        *nlp.Client
	dictionaryClient *dictionary.Client
}

func NewService(model *Model, nlpClient *nlp.Client, dictionaryClient *dictionary.Client) Service {
	return Service{
		model:            model,
		nlpClient:        nlpClient,
		dictionaryClient: dictionaryClient,
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

func (s *Service) addPage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	addPageTempl(user).Render(r.Context(), w)
}

func (s *Service) addPost(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	title := r.FormValue("title")
	content := r.FormValue("content")
	if user == nil || len(title) < 3 || len(content) < 3 {
		server.HttpError(w, http.StatusBadRequest)
		return
	}

	err := s.model.Add(Text{
		Title:   title,
		Content: content,
		UserId:  user.User.Id,
	})

	if err != nil {
		server.ServerError(w, err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Service) readPage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	textIdStr := r.PathValue("id")

	textId, err := strconv.Atoi(textIdStr)
	if err != nil {
		server.HttpError(w, http.StatusBadRequest)
		return
	}

	textPtr, err := s.model.Get(textId)
	if err != nil {
		server.ServerError(w, err)
	}
	if textPtr == nil {
		server.HttpError(w, http.StatusNotFound)
		return
	}
	text := *textPtr

	segments, err := s.splitIntoSegments(text.Content)
	if err != nil {
		server.ServerError(w, err)
	}

	readPageTempl(user, text, segments).Render(r.Context(), w)
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

	wordTempl(segment, definitions).Render(r.Context(), w)
}
