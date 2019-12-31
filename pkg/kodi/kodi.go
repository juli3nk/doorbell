package kodi

import (
	"fmt"

	"github.com/go-resty/resty"
)

type Config struct {
	URL      string
	Username string
	Password string
}

type Payload struct {
	JsonRpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      int                    `json:"id"`
}

func New(host string, port int, username, password string) (*Config, error) {
	return &Config{
		URL:      fmt.Sprintf("http://%s:%d/jsonrpc", host, port),
		Username: username,
		Password: password,
	}
}

func (c *Config) SendNotification(title, message string, displayTime int) error {
	body := Payload{
		JsonRpc: "2.0",
		Method:  "GUI.ShowNotification",
		Params:  map[string]interface{}{
			"title":       title,
			"message":     message,
			"displaytime": displayTime,
		},
		ID:     1,
	}

	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(c.Username, c.Password).
		SetBody(body).
		Post(c.URL)

	return err
}
