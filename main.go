package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/juli3nk/doorbell/pkg/kodi"
	"github.com/juli3nk/doorbell/pkg/telegram"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
)

func main() {
	var opts Options

	if err := envconfig.Process("db", &opts); err != nil {
		log.Fatal(err)
	}

	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", opts)

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

	if cc.Options.TelegramEnable {
		for _, user := range cc.Options.TelegramUsers {
			if err := telegram.Send(cc.Options.TelegramToken, user, cc.Options.TelegramMessage); err != nil {
				log.Println(err)
			}
		}
	}

	if cc.Options.KodiEnable {
		k, err := kodi.New(cc.Options.KodiHost, cc.Options.KodiPort, cc.Options.KodiUsername, cc.Options.KodiPassword)
		if err != nil {
			return err
		}

		if k.IsPlaying() {
			if err := k.SendNotification(cc.Options.KodiTitle, cc.Options.KodiMessage, cc.Options.KodiDisplayTime); err != nil {
				log.Println(err)
			}
		}
	}

	return c.NoContent(http.StatusOK)
}

func handleMute(c echo.Context) error {
	return c.String(http.StatusOK, "Muted")
}
