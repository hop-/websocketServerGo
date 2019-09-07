package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	client = http.Client{}
	// used as request pool
	semaphore = make(chan struct{}, 0)
)

type response struct {
	Status string
	Body   []byte
}

func request(method, url string,
	headers map[string]string,
	params map[string][]string,
	body interface{},
) (response, error) {
	result := response{}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return result, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return result, err
	}

	for key := range headers {
		request.Header.Set(key, headers[key])
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	qParams := request.URL.Query()
	for key := range params {
		for _, value := range params[key] {
			qParams.Add(key, value)
		}
	}

	request.URL.RawQuery = qParams.Encode()

	resp, err := client.Do(request)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	result.Body = respBody
	switch {
	case resp.StatusCode == 401:
		result.Status = "unauthorized"
	case resp.StatusCode > 400:
		result.Status = "error"
	default:
		result.Status = "ok"
	}

	return result, nil
}

func makeRequest(method, url string, headers map[string]string, params map[string][]string, body interface{}) (response, error) {
	semaphore <- struct{}{}
	defer func() { <-semaphore }()
	return request(method, url, headers, params, body)
}
