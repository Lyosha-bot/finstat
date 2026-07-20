package repository

import (
	"context"
	"finstat/internal/lib"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	INSERT_TRANSACTION_QUERY = `
		INSERT INTO transactions (user_id, value, category_id, description, date)
		SELECT $1, $2, id, $4, $5
		FROM categories
		WHERE id = $3 AND (user_id IS NULL OR user_id = $1)
		RETURNING id;
	`

	UPDATE_TRANSACTION_QUERY = `
		UPDATE transactions
		SET 
			value = $3,
			category_id = $4,
			description = $5,
			date = $6
		WHERE id = $2 AND user_id = $1
		RETURNING TRUE;
	`

	DELETE_TRANSACTION_QUERY = `
		DELETE FROM transactions
		WHERE id = $2 AND user_id = $1;
	`

	TRANSACTION_BY_ID_QUERY = `
		SELECT (
			id,
			user_id,
			value,
			description,
			date
		)
		FROM transactions 
		WHERE user_id = $1 AND id = $2;
	`

	TRANSACTIONS_QUERY_BASE = `
		SELECT (
			id,
			user_id,
			value,
			category_id,
			description,
			date
		)
		FROM transactions
		WHERE
			user_id = $1
	`
)

type Transaction struct {
	ID          uint            `json:"id" db:"id"`
	UserID      uint            `json:"userID" db:"user_id"`
	Value       decimal.Decimal `json:"value" db:"value"`
	CategoryID  uint            `json:"category_id" db:"category_id"`
	Description string          `json:"description" db:"description"`
	Date        time.Time       `json:"date" db:"date"`
}

func (c *Client) InsertTransaction(userID uint, value decimal.Decimal, categoryID uint, description string, date time.Time) (uint, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return 0, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, INSERT_TRANSACTION_QUERY, userID, value, categoryID, description, date)

	var id uint
	err = row.Scan(&id)
	if err != nil {
		return 0, lib.Ewrap("Couldn't get ID of new transaction", err)
	}

	return id, nil
}

func (c *Client) UpdateTransaction(userID uint, transactionID uint, newValue decimal.Decimal, newCategoryID uint, newDescription string, newDate time.Time) (bool, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return false, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	cmdTag, err := conn.Exec(ctx, UPDATE_TRANSACTION_QUERY, userID, transactionID, newValue, newCategoryID, newDescription, newDate)

	if err != nil {
		return false, lib.Ewrap("Couldn't update transaction", err)
	}

	return cmdTag.RowsAffected() != 0, nil
}

func (c *Client) DeleteTransaction(userID uint, transactionID uint) (bool, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return false, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	cmdTag, err := conn.Exec(ctx, DELETE_TRANSACTION_QUERY, userID, transactionID)

	if err != nil {
		return false, lib.Ewrap("Couldn't delete transaction", err)
	}

	return cmdTag.RowsAffected() != 0, nil
}

func (c *Client) TransactionByID(userID, transactionID uint) (*Transaction, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, TRANSACTION_BY_ID_QUERY, transactionID)

	var transaction Transaction
	if err = row.Scan(&transaction); err != nil {
		return nil, lib.Ewrap("Couldn't get transaction by ID", err)
	}

	return &transaction, nil
}

func (c *Client) Transactions(userID, limit, page uint, from, to *time.Time, transactionType int, categories []uint) ([]Transaction, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	args := make([]any, 0, 7)
	args = append(args, userID)

	applyArg := func(arg any) string {
		args = append(args, arg)
		return fmt.Sprintf("$%d", len(args))
	}

	var builder strings.Builder

	builder.Grow(100)

	builder.WriteString(TRANSACTIONS_QUERY_BASE)

	if from != nil {
		builder.WriteString("\n\t\tAND date >= ")
		builder.WriteString(applyArg(*from))
	}

	if to != nil {
		builder.WriteString("\n\t\tAND date <= ")
		builder.WriteString(applyArg(*to))
	}

	if transactionType < 0 {
		builder.WriteString("\n\t\tAND value < 0")
	} else if transactionType > 0 {
		builder.WriteString("\n\t\tAND value > 0")
	}

	if categories != nil {
		builder.WriteString("\n\t\tAND category_id = ANY(")
		builder.WriteString(applyArg(categories))
		builder.WriteString(")")
	}

	builder.WriteString("\nLIMIT ")
	builder.WriteString(applyArg(limit))

	offset := (page - 1) * limit

	builder.WriteString("\nOFFSET ")
	builder.WriteString(applyArg(offset))

	builder.WriteString(";")

	rows, err := conn.Query(ctx, builder.String(), args...)
	if err != nil {
		return nil, lib.Ewrap("Couldn't get transactions", err)
	}

	result := make([]Transaction, 0, limit)
	for rows.Next() {
		var val Transaction
		if err = rows.Scan(&val); err != nil {
			return nil, lib.Ewrap("Couldn't scan transaction", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, lib.Ewrap("Couldn't iterate transactions", err)
	}

	return result, nil
}
