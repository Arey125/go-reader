package dictionary

import (
	"fmt"
	"net/http"
	"reader/internal/client"
	"time"
)

type Client struct {
	httpClient http.Client
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{Timeout: time.Second * 10},
	}
}

type Definition struct {
	Definition string `json:"definition"`
}

type Meaning struct {
	Pos string `json:"partOfSpeech"`
	Definitions []Definition  `json:"definitions"`
}

type DictionaryEntry struct {
	Word string `json:"word"`
	Meanings []Meaning `json:"meanings"`
}

func (c *Client) GetEntries(word string) ([]DictionaryEntry, error) {
	url := "https://api.dictionaryapi.dev/api/v2/entries/en/" + word

	return client.MakeRequest[[]DictionaryEntry](
		&c.httpClient,
		url,
		http.MethodGet,
		nil,
	)
}
