package nlp

import (
	"net/http"
	"reader/internal/client"
	"time"
)

type Client struct {
	url        string
	httpClient http.Client
}

func NewClient(url string) Client {
	return Client{
		url:        url,
		httpClient: http.Client{Timeout: time.Second * 10},
	}
}

type Word struct {
	Text  string `json:"text"`
	Lemma string `json:"lemma"`
	Pos   string `json:"pos"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type Response struct {
	Content []Word `json:"result"`
}

func (c *Client) GetWords(text string) ([]Word, error) {
	resp, err := client.MakeRequest[Response](
		&c.httpClient,
		c.url,
		http.MethodPost,
		&client.Options{Body: map[string]string{"content": text}},
	)

	return resp.Content, err
}
