package handler

import (
	rest "app/api/sso"
	"app/bootstrap"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type RestService struct {
	dbr *pgx.Conn
	dbw *pgx.Conn
}

func RestRegister(r *echo.Echo, cfg *bootstrap.Container) {
	rest.Router(r, &RestService{
		dbr: cfg.Dbr(),
		dbw: cfg.Dbw(),
	})
}
