package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Options struct {
	Body any
}

func MakeRequest[Resp any](client *http.Client, url string, method string, options *Options) (Resp, error) {
	body := []byte("")
	if options != nil {
		reqBody, err := json.Marshal(options.Body)
		if err != nil {
			return *new(Resp), err
		}
		body = reqBody
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return *new(Resp), err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return *new(Resp), err
	}
	defer resp.Body.Close()

	respObj := *new(Resp)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respObj)
	if err != nil {
		return *new(Resp), err
	}
	return respObj, nil
}
