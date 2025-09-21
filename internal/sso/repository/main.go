package repository

import (
	database "app/database/main"

	"github.com/jackc/pgx/v5"
)

type Querier interface {
	Login
	Token
	Authorize
}

type Repository struct {
	read  Querier
	write Querier
}

func InitRepo(read *pgx.Conn, write *pgx.Conn) *Repository {
	return &Repository{
		read:  database.New(read),
		write: database.New(write),
	}
}
