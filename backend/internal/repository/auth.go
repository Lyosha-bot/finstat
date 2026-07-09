package repository

import (
	"context"
	ewrap "finstat/internal/lib"
)

const (
	INSERT_USER_QUERY = "INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id"

	USER_QUERY = "SELECT (id, name, password) FROM users WHERE name = $1"
)

type User struct {
	ID             uint
	Name           string
	HashedPassword string
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
