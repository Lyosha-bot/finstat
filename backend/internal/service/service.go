package service

import (
	"finstat/internal/models"
	"finstat/internal/repository"
	"time"

	"github.com/shopspring/decimal"
)

type Auth interface {
	Register(username, password string) error
	Login(username, password string) (accessToken string, refreshToken string, err error)
	Refresh(refreshToken string) (newAccessToken string, newRefreshToken string, err error)
	Logout(refreshToken string) error
	ID(jwtAccessToken string) (uint, error)
}

type Transaction interface {
	InsertTransaction(userID uint, value decimal.Decimal, categoryID uint, description string, date time.Time) (uint, error)
	UpdateTransaction(userID uint, transactionID uint, newValue decimal.Decimal, newCategoryID uint, newDescription string, newDate time.Time) (bool, error)
	DeleteTransaction(userID uint, transactionID uint) (bool, error)
	Transactions(userID, limit, page uint, from, to *time.Time, transactionType int, categories []uint) ([]models.Transaction, error)
}

type Category interface {
	InsertCategory(userID uint, categoryName string) (uint, error)
	UpdateCategory(userID, categoryID uint, newCategoryName string) (bool, error)
	DeleteCategory(userID, categoryID uint) (bool, error)
	SystemCategories() ([]models.Category, error)
	UserCategories(userID uint) ([]models.Category, error)
	Categories(userID uint) ([]models.Category, error)
}

type Budget interface {
	InsertBudget(userID, categoryID uint, limit decimal.Decimal) error
	UpdateBudget(userID, budgetID uint, newLimit decimal.Decimal) (bool, error)
	DeleteBudget(userID, budgetID uint) (bool, error)
	Budgets(userID uint, date time.Time) ([]models.Budget, error)
	BudgetByCategory(userID, categoryID uint, date time.Time) (*models.Budget, error)
}

type Services struct {
	Auth
	Transaction
	Category
	Budget
}

func New(repos *repository.Repository, jwtAccessSecret, jwtRefreshSecret []byte) *Services {
	return &Services{
		Auth:        NewAuthService(repos.Auth, jwtAccessSecret, jwtRefreshSecret),
		Transaction: NewTransactionService(repos.Transaction),
		Category:    NewCategoryService(repos.Category),
		Budget:      NewBudgetService(repos.Budget),
	}
}
