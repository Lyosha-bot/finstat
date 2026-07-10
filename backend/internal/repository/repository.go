// TODO: Настроить context
// TODO: Добавить интерфейсы

package repository

import (
	"context"
	"fmt"

	ewrap "finstat/internal/lib"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Credentials struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB_name  string `json:"postgres_db_name"`
}

type Client struct {
	pool *pgxpool.Pool
}

func AddClient(creds Credentials) (*Client, error) {
	ctx := context.Background()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", creds.Username, creds.Password, creds.Host, creds.Port, creds.DB_name)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't create pgxpool", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't ping database", err)
	}

	return &Client{
		pool: pool,
	}, nil
}
