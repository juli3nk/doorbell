package sound

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Config struct {
	Statefile string
	URL       string
}

func New(statefile string, host string, port int) (*Config, error) {
	return &Config{
		Statefile: statefile,
		URL:       fmt.Sprintf("http://%s:%d/play", host, port),
	}, nil
}

func (c *Config) Mute() error {
	return nil
}

func (c *Config) Unmute() error {
	return nil
}

func (c *Config) Play() error {
	client := resty.New()

	if _, err := client.R().Get(c.URL); err != nil {
		return err
	}

	return nil
}
