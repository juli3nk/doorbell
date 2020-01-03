package main

import (
	"log"
	"net/http"

	"github.com/juli3nk/doorbell/pkg/kodi"
	"github.com/juli3nk/doorbell/pkg/sound"
	"github.com/juli3nk/doorbell/pkg/telegram"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
)

func main() {
	var opts Options

	if err := envconfig.Process("doorbell", &opts); err != nil {
		log.Fatal(err)
	}

	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, &opts}
			return next(cc)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: opts.CORSAllowOrigins,
		AllowMethods: opts.CORSAllowMethods,
	}))

	e.POST("/dingdong", handleDingDong)
	e.PUT("/mute", handleMute)

	e.Logger.Fatal(e.Start(":8080"))
}

func handleDingDong(c echo.Context) error {
	cc := c.(*CustomContext)

	chT := make(chan bool)
	chK := make(chan bool)
	chS := make(chan bool)

	go func() {
		var run bool

		if cc.Options.TelegramEnable {
			for _, user := range cc.Options.TelegramUsers {
				if err := telegram.Send(cc.Options.TelegramToken, user, cc.Options.TelegramMessage); err != nil {
					c.Logger().Info(err)
				} else {
					run = true
				}
			}
		}

		chT <- run
	}()

	go func() {
		var run bool

		if cc.Options.KodiEnable {
			k, err := kodi.New(cc.Options.KodiHost, cc.Options.KodiPort, cc.Options.KodiUsername, cc.Options.KodiPassword)
			if err != nil {
				c.Logger().Info(err)
			}

			if err == nil && k.IsPlaying() {
				if err := k.SendNotification(cc.Options.KodiTitle, cc.Options.KodiMessage, cc.Options.KodiDisplayTime); err != nil {
					c.Logger().Info(err)
				} else {
					run = true
				}
			}
		}

		chK <- run
	}()

	go func() {
		var run bool

		if cc.Options.SoundEnable {
			s, err := sound.New(cc.Options.SoundStatefile, cc.Options.SoundHost, cc.Options.SoundPort)
			if err != nil {
				c.Logger().Info(err)
			}

			if err == nil {
				if err := s.Play(); err != nil {
					c.Logger().Info(err)
				} else {
					run = true
				}
			}
		}

		chS <- run
	}()

	for i := 0; i < 3; i++ {
		select {
		case r1 := <-chT:
			c.Logger().Info("Telegram notification:", r1)
		case r2 := <-chK:
			c.Logger().Info("Kodi notification:", r2)
		case r3 := <-chS:
			c.Logger().Info("Sound notification:", r3)
		}
	}

	return c.NoContent(http.StatusOK)
}

func handleMute(c echo.Context) error {
	return c.String(http.StatusOK, "Muted")
}
