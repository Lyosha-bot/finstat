package repository

import (
	"context"
	"errors"
	"finstat/internal/apperr"
	"finstat/internal/lib"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ADD_USER_QUERY = `
		INSERT INTO users (name, password)
		VALUES ($1, $2);
	`

	USER_QUERY = `
		SELECT (
			id,
			name,
			password
		)
		FROM users
		WHERE name = $1;
	`
)

type User struct {
	ID             uint
	Name           string
	HashedPassword string
}

func (c *Client) InsertUser(username, password string) error {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, ADD_USER_QUERY, username, password)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return apperr.NotUnique
			}
		}

		return lib.Ewrap("Couldn't insert user", err)
	}

	return nil
}

func (c *Client) User(username string) (*User, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, USER_QUERY, username)

	var result User
	err = row.Scan(&result)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NoRow
		}
		return nil, lib.Ewrap("Couldn't get user data", err)
	}

	return &result, nil
}
