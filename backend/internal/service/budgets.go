package service

import (
	"finstat/internal/models"
	"finstat/internal/repository"
	"time"

	"github.com/shopspring/decimal"
)

type BudgetService struct {
	repo repository.Budget
}

func getMonthPeriod(date time.Time) (from, to time.Time) {
	from = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	to = from.AddDate(0, 1, 0)
	return from, to
}

func NewBudgetService(repo repository.Budget) *BudgetService {
	return &BudgetService{
		repo: repo,
	}
}

func (s *BudgetService) InsertBudget(userID, categoryID uint, limit decimal.Decimal) error {
	return s.repo.InsertBudget(userID, categoryID, limit)
}

func (s *BudgetService) UpdateBudget(userID, budgetID uint, newLimit decimal.Decimal) (bool, error) {
	return s.repo.UpdateBudget(userID, budgetID, newLimit)
}

func (s *BudgetService) DeleteBudget(userID, budgetID uint) (bool, error) {
	return s.repo.DeleteBudget(userID, budgetID)
}

func (s *BudgetService) Budgets(userID uint, date time.Time) ([]models.Budget, error) {
	from, to := getMonthPeriod(date)
	return s.repo.Budgets(userID, from, to)
}

func (s *BudgetService) BudgetByCategory(userID, categoryID uint, date time.Time) (*models.Budget, error) {
	from, to := getMonthPeriod(date)
	return s.repo.BudgetByCategory(userID, categoryID, from, to)
}
