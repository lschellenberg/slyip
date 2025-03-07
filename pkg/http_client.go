package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	BaseUrl string
	Client  *http.Client
}

func NewHttpClient(baseUrl string) HttpClient {
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	if !strings.Contains(baseUrl, "http") {
		baseUrl = fmt.Sprintf("https://%s", baseUrl)
	}

	return HttpClient{
		BaseUrl: baseUrl,
		Client:  client,
	}
}

func (c *HttpClient) getUrl(path string) string {
	return fmt.Sprintf("%v/%v", c.BaseUrl, path)
}

func (c *HttpClient) Post(body interface{}, response interface{}, token string, path string) (int, error) {
	buffer, err := json.Marshal(body)

	if err != nil {
		return -1, err
	}

	url := c.getUrl(path)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buffer))
	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	}

	if err != nil {
		return -1, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return -1, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return -1, err
	}

	if res.StatusCode >= 300 {
		return res.StatusCode, fmt.Errorf("%s", b)
	}

	err = json.Unmarshal(b, response)

	if err != nil {
		return -1, err
	}

	return res.StatusCode, nil
}

func (c *HttpClient) Put(body interface{}, response interface{}, token string, path string) (int, error) {
	buffer, err := json.Marshal(body)

	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("PUT", c.getUrl(path), bytes.NewBuffer(buffer))

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	if err != nil {
		return -1, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return -1, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return -1, err
	}

	if res.StatusCode >= 300 {
		return res.StatusCode, fmt.Errorf("%s\n", b)
	}

	err = json.Unmarshal(b, response)

	if err != nil {
		return -1, err
	}

	return res.StatusCode, nil
}

func (c *HttpClient) Get(response interface{}, token string, path string) (int, error) {
	req, err := http.NewRequest("GET", c.getUrl(path), nil)

	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	}

	if err != nil {
		return -1, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return -1, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return -1, err
	}

	if res.StatusCode >= 300 {
		return res.StatusCode, fmt.Errorf("%s", string(b))
	}

	if response != nil {
		err = json.Unmarshal(b, response)
		if err != nil {
			return -1, err
		}
	}

	return res.StatusCode, nil
}

func (he *HttpError) read(body []byte) error {
	return json.Unmarshal(body, he)
}

type HttpError struct {
	Details string `json:"details"`
	Message string `json:"message"`
}
