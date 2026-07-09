package repository

import (
	"context"
	ewrap "finstat/internal/lib"
	"time"

	"github.com/shopspring/decimal"
)

const (
	INSERT_TRANSACTION_QUERY = `INSERT INTO transactions (user_id, amount, description, date) VALUES ($1, $2, $3, $4) RETURNING id;`

	TRANSACTION_BY_ID_QUERY = `SELECT (id, user_id, amount, description, date) FROM transactions WHERE user_id = $1 AND id = $2;`

	TRANSACTIONS_QUERY = `
		SELECT (id, user_id, amount, description, date) 
		FROM transactions 
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2
		OFFSET $3;
	`
	TRANSACTIONS_IN_PERIOD_QUERY = `
		SELECT (id, user_id, amount, description, date) 
		FROM transactions 
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
		LIMIT $4
		OFFSET $5;
	`
	TRANSACTIONS_BEFORE_DATE_QUERY = `
		SELECT (id, user_id, amount, description, date) 
		FROM transactions 
		WHERE user_id = $1 AND date <= $2
		ORDER BY date DESC
		LIMIT $3
		OFFSET $4;
	`

	TRANSACTIONS_AFTER_DATE_QUERY = `
		SELECT (id, user_id, amount, description, date) 
		FROM transactions 
		WHERE user_id = $1 AND date >= $2
		ORDER BY date DESC
		LIMIT $3
		OFFSET $4;
	`
)

type Transaction struct {
	ID          uint            `json:"id"`
	UserID      uint            `json:"userID"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
}

func (c *Client) AddTransaction(userID uint, amount decimal.Decimal, description string, date time.Time) (uint, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return 0, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, INSERT_TRANSACTION_QUERY, userID, amount, description, date)

	var id uint
	err = row.Scan(&id)
	if err != nil {
		return 0, ewrap.Wrap("Couldn't get ID of new transaction", err)
	}

	return id, nil
}

func (c *Client) TransactionByID(userID, transactionID uint) (*Transaction, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, TRANSACTION_BY_ID_QUERY, transactionID)

	var transaction Transaction
	err = row.Scan(&transaction)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get transaction by ID", err)
	}

	return &transaction, nil
}

func (c *Client) Transactions(userID, limit, page uint) ([]Transaction, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	offset := (page - 1) * limit

	rows, err := conn.Query(ctx, TRANSACTIONS_QUERY, userID, limit, offset)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get latest transactions", err)
	}

	result := make([]Transaction, 0, limit)
	for rows.Next() {
		var val Transaction
		if err = rows.Scan(&val); err != nil {
			return nil, ewrap.Wrap("Couldn't scan transaction", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, ewrap.Wrap("Couldn't iterate latest transactions", err)
	}

	return result, nil
}

func (c *Client) TransactionsInPeriod(userID uint, limit, page uint, from, to time.Time) ([]Transaction, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	offset := (page - 1) * limit

	rows, err := conn.Query(ctx, TRANSACTIONS_IN_PERIOD_QUERY, userID, from, to, limit, offset)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get transactions in period", err)
	}

	result := make([]Transaction, 0, 10)
	for rows.Next() {
		var val Transaction
		if err = rows.Scan(&val); err != nil {
			return nil, ewrap.Wrap("Couldn't scan transaction", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, ewrap.Wrap("Couldn't iterate transactions in period", err)
	}

	return result, nil
}

func (c *Client) TransactionsFromDate(userID uint, limit, page uint, date time.Time, order bool) ([]Transaction, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	offset := (page - 1) * limit

	var query string
	if order {
		query = TRANSACTIONS_BEFORE_DATE_QUERY
	} else {
		query = TRANSACTIONS_AFTER_DATE_QUERY
	}

	rows, err := conn.Query(ctx, query, userID, date, limit, offset)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get transactions before date", err)
	}

	result := make([]Transaction, 0, limit)
	for rows.Next() {
		var val Transaction
		if err = rows.Scan(&val); err != nil {
			return nil, ewrap.Wrap("Couldn't scan transaction", err)
		}

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, ewrap.Wrap("Couldn't iterate transactions from date", err)
	}

	return result, nil
}
