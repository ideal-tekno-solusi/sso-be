package bootstrap

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

func (c *Container) Dbr() *pgx.Conn {
	if c.dbr == nil {
		ctx := context.Background()

		host := viper.GetString("database.read.host")
		port := viper.GetString("database.read.port")
		database := viper.GetString("database.read.database")
		schema := viper.GetString("database.read.schema")
		username := viper.GetString("database.read.username")
		password := viper.GetString("database.read.password")
		dsn := fmt.Sprintf("host=%s user={username} password={password} dbname=%s port=%s sslmode=%s search_path=%s TimeZone=Asia/Jakarta", host, database, port, "disable", schema)

		dsn = strings.Replace(dsn, "{username}", username, 1)
		dsn = strings.Replace(dsn, "{password}", password, 1)

		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			fmt.Printf("failed to connect with error %v", err)
			return nil
		}

		c.dbr = conn
	}

	return c.dbr
}

func (c *Container) Dbw() *pgx.Conn {
	if c.dbr == nil {
		ctx := context.Background()

		host := viper.GetString("database.write.host")
		port := viper.GetString("database.write.port")
		database := viper.GetString("database.write.database")
		schema := viper.GetString("database.write.schema")
		username := viper.GetString("database.write.username")
		password := viper.GetString("database.write.password")
		dsn := fmt.Sprintf("host=%s user={username} password={password} dbname=%s port=%s sslmode=%s search_path=%s TimeZone=Asia/Jakarta", host, database, port, "disable", schema)

		dsn = strings.Replace(dsn, "{username}", username, 1)
		dsn = strings.Replace(dsn, "{password}", password, 1)

		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			fmt.Printf("failed to connect with error %v", err)
			return nil
		}

		c.dbr = conn
	}

	return c.dbr
}

func (c *Container) StopDb(ctx context.Context) {
	if c.dbr != nil {
		c.dbr.Close(ctx)
	}

	if c.dbw != nil {
		c.dbw.Close(ctx)
	}
}
