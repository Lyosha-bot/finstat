package server

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	_ "finstat/docs"

	"finstat/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	JWT_COOKIE_NAME = "refresh_jwt"
	USER_ID_KEY     = "user_id"
	REFRESH_PATH    = "/api/auth"
)

var (
	latinRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// @title Finstat API
// @version 1.0
// @description Веб-приложения для учета личных финансов
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Auth
// @in                         header
// @name                       Authorization
// @description                Укажите Access-токен

type CategoriesList []uint

func (c *CategoriesList) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	elements := strings.Split(string(text), ",")
	for _, item := range elements {
		cleaned := strings.TrimSpace(item)
		if cleaned == "" {
			continue
		}

		num, err := strconv.ParseUint(cleaned, 10, 64)
		if err != nil {
			return err
		}
		*c = append(*c, uint(num))
	}
	return nil
}

type IServer interface {
	Start(port string)
}

type Server struct {
	host     string
	services *service.Services
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type StringResponse struct {
	Result string `json:"result"`
}

func NewServer(host string, services *service.Services) *Server {
	return &Server{
		host:     host,
		services: services,
	}
}

func (s *Server) Start(port string) {
	router := gin.Default()

	api := router.Group("/api")
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := api.Group("/auth")
	auth.POST("/register", s.register)
	auth.POST("/register/is-valid", s.isValid)
	auth.POST("/login", s.login)
	auth.POST("/refresh", s.refresh)
	auth.POST("/logout", s.logout)

	transactions := api.Group("/transactions")
	transactions.Use(s.middleware)
	transactions.POST("", s.addTransaction)
	transactions.GET("", s.transactions)
	transactions.PATCH("/:id", s.updateTransaction)
	transactions.DELETE("/:id", s.deleteTransaction)

	categories := api.Group("/categories")
	categories.Use(s.middleware)
	categories.GET("", s.categories)
	categories.POST("", s.addCategory)
	categories.PATCH("/:id", s.updateCategory)
	categories.DELETE("/:id", s.deleteCategory)

	budgets := api.Group("/budgets")
	budgets.Use(s.middleware)
	budgets.POST("", s.addBudget)
	budgets.GET("", s.budgets)
	budgets.GET("/category/:id", s.budgetByCategory)
	budgets.PATCH("/:id", s.updateBudget)
	budgets.DELETE("/:id", s.deleteBudget)

	router.Run(":" + port)
}

func (s *Server) middleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		log.Println("Отсутствует header: Authorization")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует header: Authorization"})
		return
	}

	id, err := s.services.Auth.ID(token)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}

	c.Set(USER_ID_KEY, id)

	c.Next()
}
