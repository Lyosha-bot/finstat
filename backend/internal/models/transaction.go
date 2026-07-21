package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID          uint            `json:"id" db:"id"`
	UserID      uint            `json:"userID" db:"user_id"`
	Value       decimal.Decimal `json:"value" db:"value"`
	CategoryID  uint            `json:"category_id" db:"category_id"`
	Description string          `json:"description" db:"description"`
	Date        time.Time       `json:"date" db:"date"`
}
