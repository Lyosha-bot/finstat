package repository

import (
	"context"
	"errors"
	"finstat/internal/apperr"
	"finstat/internal/lib"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const ACTIVE_REFRESH_TOKENS_COUNT = 4

const (
	INSERT_USER_QUERY = `
		INSERT INTO users (name, password)
		VALUES ($1, $2);
	`

	INSERT_REFRESH_TOKEN_QUERY = `
		WITH inserted AS (
			INSERT INTO refresh_tokens (user_id, expires_at)
			VALUES ($1, $2)
			RETURNING uuid
		),
		deleted AS (
			DELETE FROM refresh_tokens
			WHERE user_id = $1
			AND uuid NOT IN (
				SELECT uuid 
				FROM refresh_tokens 
				WHERE user_id = $1
				ORDER BY created_at DESC
				LIMIT $3
			)
		)
		SELECT uuid FROM inserted;
	`

	DELETE_REFRESH_TOKEN_QUERY = `
		DELETE FROM refresh_tokens
		WHERE uuid = $1;
	`

	DELETE_ALL_REFRESH_TOKENS_QUERY = `
		DELETE FROM refresh_tokens
		WHERE user_id = $1;
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

	REFRESH_TOKEN_QUERY = `
		SELECT (
			uuid,
			user_id,
			expires_at
		)
		FROM refresh_tokens
		WHERE uuid = $1;
	`
)

type User struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	HashedPassword string `json:"hashed_password"`
}

type RefreshToken struct {
	UUID      string    `json:"uuid"`
	UserID    uint      `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (c *Client) InsertUser(username, password string) error {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, INSERT_USER_QUERY, username, password)

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

func (c *Client) InsertRefreshToken(userID uint, expiresAt time.Time) (string, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return "", lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, INSERT_REFRESH_TOKEN_QUERY, userID, expiresAt, ACTIVE_REFRESH_TOKENS_COUNT)

	var id string
	if err := row.Scan(&id); err != nil {
		return "", lib.Ewrap("Couldn't insert refresh token", err)
	}

	return id, nil
}

func (c *Client) DeleteRefreshToken(tokenUUID string) (bool, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return false, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	cmdTag, err := conn.Exec(ctx, DELETE_REFRESH_TOKEN_QUERY, tokenUUID)

	if err != nil {
		return false, lib.Ewrap("Couldn't delete refresh token", err)
	}

	return cmdTag.RowsAffected() != 0, nil
}

func (c *Client) DeleteAllRefreshTokens(userID uint) (bool, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return false, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	cmdTag, err := conn.Exec(ctx, DELETE_ALL_REFRESH_TOKENS_QUERY, userID)

	if err != nil {
		return false, lib.Ewrap("Couldn't delete all refresh tokens", err)
	}

	return cmdTag.RowsAffected() != 0, nil
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
	if err := row.Scan(&result); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NoRows
		}
		return nil, lib.Ewrap("Couldn't get user data", err)
	}

	return &result, nil
}

func (c *Client) RefreshToken(uuid string) (*RefreshToken, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, REFRESH_TOKEN_QUERY, uuid)

	var result RefreshToken
	if err := row.Scan(&result); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NoRows
		}
		return nil, lib.Ewrap("Couldn't get refresh token data", err)
	}

	return &result, nil
}
