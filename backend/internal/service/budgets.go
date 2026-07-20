package service

import (
	"finstat/internal/repository"
	"time"

	"github.com/shopspring/decimal"
)

type Budget = repository.Budget

type BudgetRepo interface {
	AddBudget(userID, categoryID uint, limit decimal.Decimal) error
	UpdateBudget(userID, budgetID uint, newLimit decimal.Decimal) (bool, error)
	DeleteBudget(userID, budgetID uint) (bool, error)
	Budgets(userID uint, from, to time.Time) ([]Budget, error)
	BudgetByCategory(userID, categoryID uint, from, to time.Time) (*Budget, error)
}

type BudgetService struct {
	repo BudgetRepo
}

func getMonthPeriod(date time.Time) (from, to time.Time) {
	from = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	to = from.AddDate(0, 1, 0)
	return from, to
}

func NewBudgetService(repo BudgetRepo) *BudgetService {
	return &BudgetService{
		repo: repo,
	}
}

func (s *BudgetService) AddBudget(userID, categoryID uint, limit decimal.Decimal) error {
	return s.repo.AddBudget(userID, categoryID, limit)
}

func (s *BudgetService) UpdateBudget(userID, budgetID uint, newLimit decimal.Decimal) (bool, error) {
	return s.repo.UpdateBudget(userID, budgetID, newLimit)
}

func (s *BudgetService) DeleteBudget(userID, budgetID uint) (bool, error) {
	return s.repo.DeleteBudget(userID, budgetID)
}

func (s *BudgetService) Budgets(userID uint, date time.Time) ([]Budget, error) {
	from, to := getMonthPeriod(date)
	return s.repo.Budgets(userID, from, to)
}

func (s *BudgetService) BudgetByCategory(userID, categoryID uint, date time.Time) (*Budget, error) {
	from, to := getMonthPeriod(date)
	return s.repo.BudgetByCategory(userID, categoryID, from, to)
}
