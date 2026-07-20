package repository

import (
	"context"
	"finstat/internal/lib"
)

const (
	INSERT_CATEGORY_QUERY = `
		INSERT INTO categories (user_id, name) 
		VALUES ($1, $2) 
		RETURNING id
	`

	UPDATE_CATEGORY_QUERY = `
		UPDATE categories
		SET 
			name = $3
		WHERE id = $2 AND user_id = $1;
	`

	DELETE_CATEGORY_QUERY = `
		DELETE FROM categories
		WHERE id = $2 AND user_id = $1;
	`

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

func (c *Client) InsertCategory(userID uint, categoryName string) (uint, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return 0, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, INSERT_CATEGORY_QUERY, userID, categoryName)

	var id uint
	err = row.Scan(&id)
	if err != nil {
		return 0, lib.Ewrap("Couldn't get ID of new category", err)
	}

	return id, nil
}

func (c *Client) UpdateCategory(userID, categoryID uint, newCategoryName string) (bool, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return false, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	cmdTag, err := conn.Exec(ctx, UPDATE_CATEGORY_QUERY, userID, categoryID, newCategoryName)

	if err != nil {
		return false, lib.Ewrap("Couldn't update category", err)
	}

	return cmdTag.RowsAffected() != 0, nil
}

func (c *Client) DeleteCategory(userID, categoryID uint) (bool, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return false, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	cmdTag, err := conn.Exec(ctx, DELETE_CATEGORY_QUERY, userID, categoryID)

	if err != nil {
		return false, lib.Ewrap("Couldn't delete category", err)
	}

	return cmdTag.RowsAffected() != 0, nil
}

func (c *Client) SystemCategories() ([]Category, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, SYSTEM_CATEGORIES_QUERY)
	if err != nil {
		return nil, lib.Ewrap("Couldn't get system categories", err)
	}

	result := make([]Category, 0, 10)
	for rows.Next() {
		var val Category
		if err = rows.Scan(&val); err != nil {
			return nil, lib.Ewrap("Couldn't scan category", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, lib.Ewrap("Couldn't iterate system categories", err)
	}

	return result, nil
}

func (c *Client) UserCategories(userID uint) ([]Category, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, USER_CATEGORIES_QUERY, userID)
	if err != nil {
		return nil, lib.Ewrap("Couldn't get user categories", err)
	}

	result := make([]Category, 0, 10)
	for rows.Next() {
		var val Category
		if err = rows.Scan(&val); err != nil {
			return nil, lib.Ewrap("Couldn't scan category", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, lib.Ewrap("Couldn't iterate user categories", err)
	}

	return result, nil
}

func (c *Client) Categories(userID uint) ([]Category, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, CATEGORIES_QUERY, userID)
	if err != nil {
		return nil, lib.Ewrap("Couldn't get categories", err)
	}

	result := make([]Category, 0, 10)
	for rows.Next() {
		var val Category
		if err = rows.Scan(&val); err != nil {
			return nil, lib.Ewrap("Couldn't scan category", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, lib.Ewrap("Couldn't iterate categories", err)
	}

	return result, nil
}
