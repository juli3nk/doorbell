package main

import (
	"github.com/labstack/echo"
)

type Options struct {
	CORSAllowOrigins []string `split_words:"true" default:"*"`
	CORSAllowMethods []string `split_words:"true" default:"POST,PUT"`
	SoundEnable      bool     `split_words:"true" default:"false"`
	SoundStatefile   string   `split_words:"true" default:"/tmp/doorbell.state"`
	SoundHost        string   `split_words:"true"`
	SoundPort        int      `split_words:"true"`
	TelegramEnable   bool     `split_words:"true" default:"false"`
	TelegramToken    string   `split_words:"true"`
	TelegramUsers    []int    `split_words:"true"`
	TelegramMessage  string   `split_words:"true" default:"Ding Dong"`
	KodiEnable       bool     `split_words:"true" default:"false"`
	KodiHost         string   `split_words:"true"`
	KodiPort         int      `split_words:"true"`
	KodiUsername     string   `split_words:"true"`
	KodiPassword     string   `split_words:"true"`
	KodiTitle        string   `split_words:"true" default:"DoorBell"`
	KodiMessage      string   `split_words:"true" default:"Ding Dong"`
	KodiDisplayTime  int      `split_words:"true" default:"5000"`
}

type CustomContext struct {
	echo.Context

	Options *Options
}
