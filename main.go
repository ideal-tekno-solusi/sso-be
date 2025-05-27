package main

import (
	"app/api"
	"app/bootstrap"
	"app/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	vd "app/api/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	r := echo.New()

	r.Logger.SetLevel(log.INFO)
	r.Validator = &vd.CustomValidator{Validator: validator.New()}

	csrfDomain := viper.GetString("config.csrf.domain")
	csrfPath := viper.GetString("config.csrf.path")
	csrfAge := viper.GetInt("config.csrf.age")

	// TODO: cek lagi CORS ini
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           36000,
	}))

	r.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     csrfPath,
		CookieDomain:   csrfDomain,
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookieMaxAge:   csrfAge,
		TokenLookup:    "cookie:_csrf",
		ErrorHandler: func(err error, c echo.Context) error {
			message := err.(*echo.HTTPError)
			utils.SendProblemDetailJson(c, message.Code, message.Message.(string), c.Path(), uuid.NewString())

			return nil
		},
	}))

	cfg := bootstrap.InitContainer()

	api.RegisterApi(r, cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := r.Start(fmt.Sprintf("%v:%v", viper.GetString("services.host"), viper.GetString("services.port"))); err != nil && err != http.ErrServerClosed {
			r.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg.StopDb(ctx)

	if err := r.Shutdown(ctx); err != nil {
		r.Logger.Fatal(err)
	}

	logrus.Warn("Server exiting")
}
