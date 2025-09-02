package nlp

import (
	"bytes"
	"encoding/json"
	"net/http"
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
	body, err := json.Marshal(map[string]string{"content": text})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respObj := Response{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respObj)
	if err != nil {
		return nil, err
	}
	return respObj.Content, nil
}
