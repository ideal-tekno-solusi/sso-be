package main

import (
	"app/api"
	"app/bootstrap"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	vd "app/api/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

	cfg := bootstrap.InitContainer()

	allowOrigins := viper.GetStringSlice("config.cors.allow_origins")

	// TODO: cek lagi CORS ini
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		ExposeHeaders:    []string{echo.HeaderContentType, echo.HeaderContentLength},
		AllowCredentials: true,
		MaxAge:           36000,
	}))

	//? set up config for session
	key := viper.GetString("config.session.secret")
	domain := viper.GetString("config.session.domain")
	store := sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		Domain:   domain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	r.Use(session.Middleware(store))

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
