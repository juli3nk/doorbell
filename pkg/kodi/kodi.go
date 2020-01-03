package kodi

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Config struct {
	URL      string
	Username string
	Password string
}

type Payload struct {
	JsonRpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
	ID      int                    `json:"id"`
}

type Result struct {
	JsonRpc string                   `json:"jsonrpc"`
	ID      int                      `json:"id"`
	Result  []map[string]interface{} `json:"result,omitempty"`
}

func New(host string, port int, username, password string) (*Config, error) {
	return &Config{
		URL:      fmt.Sprintf("http://%s:%d/jsonrpc", host, port),
		Username: username,
		Password: password,
	}, nil
}

func (c *Config) IsPlaying() bool {
	body := Payload{
		JsonRpc: "2.0",
		Method:  "Player.GetActivePlayers",
		ID:      1,
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(c.Username, c.Password).
		SetBody(body).
		SetResult(&Result{}).
		Post(c.URL)

	if err != nil {
		return false
	}

	result := resp.Result().(*Result)
	if len(result.Result) > 0 {
		return true
	}

	return false
}

func (c *Config) SendNotification(title, message string, displayTime int) error {
	body := Payload{
		JsonRpc: "2.0",
		Method:  "GUI.ShowNotification",
		Params: map[string]interface{}{
			"title":       title,
			"message":     message,
			"displaytime": displayTime,
		},
		ID: 1,
	}

	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(c.Username, c.Password).
		SetBody(body).
		Post(c.URL)

	return err
}
