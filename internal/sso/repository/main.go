package repository

import (
	database "app/database/main"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	read  *database.Queries
	write *database.Queries
}

func InitRepo(read *pgx.Conn, write *pgx.Conn) *Repository {
	return &Repository{
		read:  database.New(read),
		write: database.New(write),
	}
}
