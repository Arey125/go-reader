package reader

import (
	"net/http"
	"reader/internal/dictionary"
	"reader/internal/nlp"
	"reader/internal/server"
	"reader/internal/users"
)

type Service struct {
	textModel        *TextModel
	wordModel        *WordModel
	nlpClient        *nlp.Client
	dictionaryClient *dictionary.Client
	wordFreq         WordFreq
}

func NewService(textModel *TextModel, wordModel *WordModel, nlpClient *nlp.Client, dictionaryClient *dictionary.Client) Service {
	return Service{
		textModel:        textModel,
		wordModel:        wordModel,
		nlpClient:        nlpClient,
		dictionaryClient: dictionaryClient,
		wordFreq:         NewWordFreq(),
	}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", s.homePage)
	mux.HandleFunc("GET /texts/add", s.addPage)
	mux.HandleFunc("POST /texts/add", s.addPost)
	mux.HandleFunc("GET /texts/{id}", s.readPage)
	mux.HandleFunc("POST /texts/{id}/save-words", s.saveWordsPost)

	mux.HandleFunc("GET /word", s.wordGet)
}

func (s *Service) homePage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	texts, err := s.textModel.All()
	if err != nil {
		server.ServerError(w, err)
	}
	homePageTempl(user, texts).Render(r.Context(), w)
}

func (s *Service) readPage(w http.ResponseWriter, r *http.Request) {
	user := users.GetUser(r)
	pagePtr, err := s.getPage(w, r.PathValue("id"), r.FormValue("page"))
	if err != nil {
		return
	}

	if r.FormValue("cur") != "" {
		curPagePtr, err := s.getPage(w, r.PathValue("id"), r.FormValue("cur"))
		if err != nil {
			server.HttpError(w, http.StatusBadRequest)
		}
		if curPagePtr != nil {
			s.saveWordsFromPage(*curPagePtr)
		}
	}

	page := *pagePtr
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

func (s *Service) saveWordsPost(w http.ResponseWriter, r *http.Request) {
	pagePtr, err := s.getPage(w, r.PathValue("id"), r.FormValue("page"))
	if err != nil {
		return
	}

	page := *pagePtr
	err = s.saveWordsFromPage(page)
	if err != nil {
		server.ServerError(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
