package repository

import (
	"context"
	ewrap "finstat/internal/lib"
)

const (
	ADD_CATEGORY_QUERY = `INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id`

	SYSTEM_CATEGORIES_QUERY = `
		SELECT (
			id,
			name
		)
		FROM categories
		WHERE user_id is NULL;
	`

	USER_CATEGORIES_QUERY = `
		SELECT (
			id,
			name
		)
		FROM categories
		WHERE user_id = $1;
	`

	CATEGORIES_QUERY = `
		SELECT (
			id,
			name
		)
		FROM categories
		WHERE user_id IS NULL OR user_id = $1;
	`
)

type Category struct {
	ID   uint   `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (c *Client) AddCategory(userID uint, categoryName string) (uint, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return 0, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, ADD_CATEGORY_QUERY, userID, categoryName)

	var id uint
	err = row.Scan(&id)
	if err != nil {
		return 0, ewrap.Wrap("Couldn't get ID of new category", err)
	}

	return id, nil
}

func (c *Client) SystemCategories() ([]Category, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, SYSTEM_CATEGORIES_QUERY)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get system categories", err)
	}

	result := make([]Category, 0, 10)
	for rows.Next() {
		var val Category
		if err = rows.Scan(&val); err != nil {
			return nil, ewrap.Wrap("Couldn't scan category", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, ewrap.Wrap("Couldn't iterate system categories", err)
	}

	return result, nil
}

func (c *Client) UserCategories(userID uint) ([]Category, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, USER_CATEGORIES_QUERY, userID)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get user categories", err)
	}

	result := make([]Category, 0, 10)
	for rows.Next() {
		var val Category
		if err = rows.Scan(&val); err != nil {
			return nil, ewrap.Wrap("Couldn't scan category", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, ewrap.Wrap("Couldn't iterate user categories", err)
	}

	return result, nil
}

func (c *Client) Categories(userID uint) ([]Category, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, CATEGORIES_QUERY, userID)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get categories", err)
	}

	result := make([]Category, 0, 10)
	for rows.Next() {
		var val Category
		if err = rows.Scan(&val); err != nil {
			return nil, ewrap.Wrap("Couldn't scan category", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, ewrap.Wrap("Couldn't iterate categories", err)
	}

	return result, nil
}
