package repository

import (
	"context"
	"errors"
	"finstat/internal/apperr"
	"finstat/internal/lib"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
)

const (
	ADD_BUDGET_QUERY = `
		INSERT INTO budgets (user_id, category_id, limit_value) 
		VALUES ($1, $2, $3);
	`

	UPDATE_BUDGET_QUERY = `
		UPDATE budgets
		SET 
			limit_value = $3
		WHERE id = $2 AND user_id = $1;
	`

	DELETE_BUDGET_QUERY = `
		DELETE FROM budgets
		WHERE id = $2 AND user_id = $1;
	`

	BUDGETS_QUERY = `
		SELECT 
			b.id,
			c.id AS category_id,
			COALESCE(c.name, 'Все категории') AS category_name,
			b.limit_value,
			COALESCE(SUM(t.value), 0)
		FROM budgets b
		LEFT JOIN categories c ON c.id = b.category_id
		LEFT JOIN transactions t ON t.user_id = b.user_id 
			AND (b.category_id IS NULL OR t.category_id = b.category_id) 
			AND t.date >= $2 
			AND t.date < $3
			AND t.value < 0
		WHERE b.user_id = $1
		GROUP BY b.id, c.id;
	`

	BUDGET_BY_CATEGORY_QUERY = `
		SELECT
			b.id,
			c.id AS category_id,
			COALESCE(c.name, 'Все категории') AS category_name,
			b.limit_value,
			COALESCE(SUM(t.value), 0) AS current_value
		FROM budgets b
		LEFT JOIN categories c ON c.id = b.category_id
		LEFT JOIN transactions t ON t.user_id = b.user_id 
			AND (b.category_id IS NULL OR t.category_id = b.category_id) 
			AND t.date >= $3 
			AND t.date < $4
			AND t.value < 0
		WHERE b.user_id = $1 AND b.category_id = $2
		GROUP BY b.id, c.id;
	`
)

type Budget struct {
	ID           uint            `json:"id" db:"id"`
	CategoryID   uint            `json:"category_id" db:"category_id"`
	CategoryName string          `json:"category_name" db:"category_name"`
	LimitValue   decimal.Decimal `json:"limit_value" db:"limit_value"`
	CurrentValue decimal.Decimal `json:"current_value" db:"current_value"`
}

func (c *Client) AddBudget(userID, categoryID uint, limit decimal.Decimal) error {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, ADD_BUDGET_QUERY, userID, categoryID, limit)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return apperr.NotUnique
			}
		}

		return lib.Ewrap("Couldn't add new budget", err)
	}

	return nil
}

func (c *Client) UpdateBudget(userID, budgetID uint, newLimit decimal.Decimal) (bool, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return false, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	cmdTag, err := conn.Exec(ctx, UPDATE_BUDGET_QUERY, userID, budgetID, newLimit)

	if err != nil {
		return false, lib.Ewrap("Couldn't update budget", err)
	}

	return cmdTag.RowsAffected() != 0, nil
}

func (c *Client) DeleteBudget(userID, budgetID uint) (bool, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return false, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	cmdTag, err := conn.Exec(ctx, DELETE_BUDGET_QUERY, userID, budgetID)

	if err != nil {
		return false, lib.Ewrap("Couldn't delete budget", err)
	}

	return cmdTag.RowsAffected() != 0, nil
}

func (c *Client) Budgets(userID uint, from, to time.Time) ([]Budget, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, BUDGETS_QUERY, userID, from, to)
	if err != nil {
		return nil, lib.Ewrap("Couldn't get budgets data", err)
	}

	result := make([]Budget, 0, 5)
	for rows.Next() {
		var val Budget
		if err = rows.Scan(&val.ID, &val.CategoryID, &val.CategoryName, &val.LimitValue, &val.CurrentValue); err != nil {
			return nil, lib.Ewrap("Couldn't scan budget", err)
		}

		val.CurrentValue = val.CurrentValue.Abs()

		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return nil, lib.Ewrap("Couldn't iterate budgets", err)
	}

	return result, nil
}

func (c *Client) BudgetByCategory(userID, categoryID uint, from, to time.Time) (*Budget, error) {
	ctx := context.Background()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't acquire connection", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, BUDGET_BY_CATEGORY_QUERY, userID, categoryID, from, to)

	var result Budget
	err = row.Scan(&result.ID, &result.CategoryID, &result.CategoryName, &result.LimitValue, &result.CurrentValue)
	if err != nil {
		return nil, lib.Ewrap("Couldn't get budget data", err)
	}

	result.CurrentValue = result.CurrentValue.Abs()

	return &result, nil
}
