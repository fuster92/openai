package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/philippta/trip"
	"net/http"
	"time"
)

const baseUrl = "https://api.openai.com"

type Client struct {
	httpClient *http.Client
	BaseUrl    string
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) Chat(r *ChatRequest) (*ChatResponse, error) {
	chatEndpoint := "chat/completions"
	data, err := json.Marshal(&r)
	if err != nil {
		return nil, fmt.Errorf("error sending chat request: %v", err)
	}
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/%s", c.BaseUrl, chatEndpoint),
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, fmt.Errorf("error reading chat response: %v", err)
	}
	defer resp.Body.Close()
	response := ChatResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding chat response: %v", err)
	}
	return &response, nil
}

func NewClient(token, version string) *Client {
	var (
		attempts = 3
		delay    = 50 * time.Millisecond
	)

	client := &http.Client{
		Transport: trip.Default(
			trip.BearerToken(token),
			trip.Retry(attempts, delay, trip.RetryableStatusCodes...),
		),
	}
	return &Client{
		httpClient: client,
		BaseUrl:    fmt.Sprintf("%s/%s", baseUrl, version),
	}
}
