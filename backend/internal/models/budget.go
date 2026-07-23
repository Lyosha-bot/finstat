package models

import "github.com/shopspring/decimal"

type Budget struct {
	ID           uint            `json:"id" db:"id"`
	CategoryID   uint            `json:"category_id" db:"category_id"`
	CategoryName string          `json:"category_name" db:"category_name"`
	LimitValue   decimal.Decimal `json:"limit_value" db:"limit_value"`
	CurrentValue decimal.Decimal `json:"current_value" db:"current_value"`
}
