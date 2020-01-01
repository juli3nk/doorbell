package main

import (
	"github.com/labstack/echo"
)

type Options struct {
	CORSAllowOrigins []string `split_words:"true" default:"*"`
	CORSAllowMethods []string `split_words:"true" default:"POST"`
	TelegramEnable   bool     `split_words:"true" default:"false"`
	TelegramToken    string   `split_words:"true"`
	TelegramUsers    []int    `split_words:"true"`
	TelegramMessage  string   `split_words:"true" default:"Ding Dong"`
	KodiEnable       bool     `split_words:"true" default:"false"`
	KodiHost         string   `split_words:"true" validate:"omitempty,tcp4_addr"`
	KodiPort         int      `split_words:"true"`
	KodiUsername     string   `split_words:"true"`
	KodiPassword     string   `split_words:"true"`
	KodiTitle        string   `split_words:"true" default:"DoorBell"`
	KodiMessage      string   `split_words:"true" default:"Ding Dong"`
	KodiDisplayTime  int      `split_words:"true" default:"5000"`
	SoundEnable      bool     `split_words:"true" default:"false"`
}

type CustomContext struct {
	echo.Context

	Options *Options
}
