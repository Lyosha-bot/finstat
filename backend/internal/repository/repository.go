// TODO: Настроить context
// TODO: Добавить интерфейсы

package repository

import (
	"context"
	"fmt"

	ewrap "auth.my-financials/internal/lib"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	INSERT_USER_QUERY = "INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id"

	USER_QUERY = "SELECT (id, name, password) FROM users WHERE name = $1"
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

type User struct {
	ID             uint
	Name           string
	HashedPassword string
}

func NewClient(creds Credentials) (*Client, error) {
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

func (c *Client) InsertUser(username, password string) (uint, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return 0, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, INSERT_USER_QUERY, username, password)

	var id uint
	err = row.Scan(&id)
	if err != nil {
		return 0, ewrap.Wrap("Couldn't get ID of new user", err)
	}

	return id, nil
}

func (c *Client) User(username string) (*User, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, USER_QUERY, username)

	var result User
	err = row.Scan(&result)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get user data", err)
	}

	return &result, nil
}
