package server

import (
	"log"
	"net/http"
	"time"

	_ "finstat/docs"

	"finstat/internal/repository"
	"finstat/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	TOKEN_LIFE_TIME = 20 // В минутах
	JWT_COOKIE_NAME = "jwt"
)

// @title Finstat API
// @version 0.5
// @description Веб-приложения для учета личных финансов
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api

type IServer interface {
	Start(port string)
}

type Server struct {
	host                string
	authService         *service.AuthService
	transactionsService *service.TransactionsService
}

type UserFormat struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}

type AddTransactionFormat struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Description string          `json:"description" binding:"omitempty"`
	Date        string          `json:"date" binding:"required"`
}

type TransactionsFilter struct {
	Limit uint   `form:"limit" binding:"numeric,min=1"`
	Page  uint   `form:"page" binding:"numeric,min=1"`
	From  string `form:"from" binding:"omitempty,datetime=2006-01-02"`
	To    string `form:"to" binding:"omitempty,datetime=2006-01-02"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TransactionsResponse struct {
	Result []repository.Transaction `json:"result"`
}

func NewServer(host string, authService *service.AuthService, transactionsService *service.TransactionsService) *Server {
	return &Server{
		host:                host,
		authService:         authService,
		transactionsService: transactionsService,
	}
}

func (s *Server) Start(port string) {
	router := gin.Default()

	api := router.Group("/api")
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := api.Group("/auth")
	auth.POST("/register", s.register)
	auth.POST("/login", s.login)

	transactions := api.Group("/transactions")
	transactions.Use(s.middleware)
	transactions.POST("", s.addTransaction)
	transactions.GET("", s.transactions)

	router.Run(":" + port)
}

func (s *Server) middleware(c *gin.Context) {
	cookie, err := c.Cookie(JWT_COOKIE_NAME)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		c.Abort()
		return
	}

	id, err := s.authService.ID(cookie)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		c.Abort()
		return
	}

	c.Set("jwt", id)

	c.Next()
}

// @Summary 		Регистрация нового пользователя
// @Description  	Создает аккаунт в системе.
// @Tags         	auth
// @Accept       	json
// @Produce      	json
// @Param        	input body UserFormat true 	"Логин и пароль для регистрации"
// @Success      	201  {object}  MessageResponse 	"Успешная регистрация"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка регистрации пользователя"
// @Router       	/auth/register [post]
func (s *Server) register(c *gin.Context) {
	var data UserFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	if err := s.authService.Register(data.Username, data.Password); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Успешная регистрация"})
}

// @Summary 		Авторизация пользователя
// @Description  	Авторизирует пользователя системы и сохраняет jwt-токена в cookies
// @Tags         	auth
// @Accept       	json
// @Produce      	json
// @Param        	input body UserFormat true 		"Логин и пароль для авторизации"
// @Success      	200  {object}  MessageResponse 	"Успешная авторизация"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка авторизации пользователя"
// @Router       	/auth/login [post]
func (s *Server) login(c *gin.Context) {
	var data UserFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	newToken, err := s.authService.Login(data.Username, data.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка авторизации пользователя"})
		return
	}

	c.SetCookie(
		JWT_COOKIE_NAME,
		newToken,
		TOKEN_LIFE_TIME*60,
		"/",
		s.host,
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Успешная авторизация"})
}

// @Summary 		Создание транзакции
// @Description  	Создает транзакцию от имени авторизированного пользователя
// @Tags         	transactions
// @Accept       	json
// @Produce      	json
// @Param        	input body AddTransactionFormat true "Информация о транзакции"
// @Success      	200  {object}  MessageResponse 	"Успешно"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	401  {object}  ErrorResponse 	"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при добавления транзакции"
// @Router       	/transactions [post]
func (s *Server) addTransaction(c *gin.Context) {
	var data AddTransactionFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	date, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнена дата"})
		return
	}

	id := c.GetUint("jwt")

	_, err = s.transactionsService.AddTransaction(id, data.Amount, data.Description, date)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавления транзакции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешно"})
}

// @Summary 		Получение списка транзакций
// @Description  	Возвращает список транзакций выполненных авторизованным пользователем
// @Tags         	transactions
// @Produce      	json
// @Param        	query query TransactionsFilter true "Параметры полученных транзакций"
// @Success      	200  {object}  TransactionsResponse 	"Успешное получение транзакций"
// @Failure      	400  {object}  ErrorResponse 			"Неверно заполнены поля"
// @Failure      	400  {object}  ErrorResponse 			"Неверно заполнена дата начала периода"
// @Failure      	401  {object}  ErrorResponse 			"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 			"Ошибка при получении транзакций"
// @Router       	/transactions [get]
func (s *Server) transactions(c *gin.Context) {
	var data TransactionsFilter
	err := c.ShouldBindQuery(&data)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	var from time.Time
	if data.From != "" {
		from, err = time.Parse("2006-01-02", data.From)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнена дата начала периода"})
			return
		}
	}

	var to time.Time
	if data.To != "" {
		to, err = time.Parse("2006-01-02", data.To)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнена дата конца периода"})
			return
		}
	}

	id := c.GetUint("jwt")

	var result []repository.Transaction

	if data.From == "" && data.To == "" {
		result, err = s.transactionsService.Transactions(id, data.Limit, data.Page)
	} else if data.To == "" {
		result, err = s.transactionsService.TransactionsFromDate(id, data.Limit, data.Page, from, true)
	} else if data.From == "" {
		result, err = s.transactionsService.TransactionsFromDate(id, data.Limit, data.Page, to, false)
	} else {
		result, err = s.transactionsService.TransactionsInPeriod(id, data.Limit, data.Page, from, to)
	}

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении транзакций"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
