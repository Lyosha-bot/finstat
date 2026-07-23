// TODO: Настроить context
// TODO: Добавить интерфейсы

package repository

import (
	"context"
	"finstat/internal/lib"
	"finstat/internal/models"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Credentials struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB_name  string `json:"postgres_db_name"`
}

type Auth interface {
	InsertUser(username, password string) error
	InsertRefreshToken(userID uint, expires_at time.Time) (string, error)
	DeleteRefreshToken(tokenUUID string) (bool, error)
	DeleteAllRefreshTokens(userID uint) (bool, error)
	User(username string) (*models.User, error)
	RefreshToken(tokenUUID string) (*models.RefreshToken, error)
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
	Budgets(userID uint, from, to time.Time) ([]models.Budget, error)
	BudgetByCategory(userID, categoryID uint, from, to time.Time) (*models.Budget, error)
}

type Repository struct {
	Auth
	Transaction
	Category
	Budget
}

func New(creds Credentials) (*Repository, error) {
	ctx := context.Background()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", creds.Username, creds.Password, creds.Host, creds.Port, creds.DB_name)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, lib.Ewrap("Couldn't create pgxpool", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, lib.Ewrap("Couldn't ping database", err)
	}

	return &Repository{
		Auth:        NewAuthRepo(pool),
		Transaction: NewTransactionRepo(pool),
		Category:    NewCategoryRepo(pool),
		Budget:      NewBudgetRepo(pool),
	}, nil
}
