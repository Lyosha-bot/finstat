package server

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "finstat/docs"

	"finstat/internal/apperr"
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

var (
	latinRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// @title Finstat API
// @version 0.7
// @description Веб-приложения для учета личных финансов
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api

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
	host               string
	authService        *service.AuthService
	transactionService *service.TransactionService
	categoryService    *service.CategoryService
	budgetService      *service.BudgetService
}

type UserFormat struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type AddTransactionFormat struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Category    uint            `json:"category" binding:"required,numeric"`
	Description string          `json:"description" binding:"omitempty"`
	Date        string          `json:"date" binding:"required"`
}

type UpdateTransactionFormat struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	CategoryID  uint            `json:"category" binding:"required,numeric"`
	Description string          `json:"description" binding:"omitempty"`
	Date        string          `json:"date" binding:"required"`
}

type DeleteTransactionFormat struct {
	ID uint `json:"id" binding:"required, numeric"`
}

type TransactionsFilter struct {
	Limit      uint           `form:"limit" binding:"numeric,min=1"`
	Page       uint           `form:"page" binding:"numeric,min=1"`
	From       string         `form:"from" binding:"omitempty,datetime=2006-01-02"`
	To         string         `form:"to" binding:"omitempty,datetime=2006-01-02"`
	Type       int            `form:"type" binding:"omitempty,numeric,min=-1,max=1"`
	Categories CategoriesList `form:"categories,parser=encoding.TextUnmarshaler" binding:"omitempty"`
}

type AddBudgetFormat struct {
	CategoryID uint            `json:"category_id" binding:"required"`
	Limit      decimal.Decimal `json:"limit" binding:"omitempty"`
}

type BudgetsFilter struct {
	Date string `form:"date" binding:"required,datetime=2006-01-02"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TransactionsResponse struct {
	Result []service.Transaction `json:"result"`
}

type CategoriesResponse struct {
	Result []service.Category `json:"result"`
}

type BudgetsResponse struct {
	Result []service.Budget `json:"result"`
}

func AddServer(host string, authService *service.AuthService, transactionsService *service.TransactionService, categoryService *service.CategoryService, budgetService *service.BudgetService) *Server {
	return &Server{
		host:               host,
		authService:        authService,
		transactionService: transactionsService,
		categoryService:    categoryService,
		budgetService:      budgetService,
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

	transactions := api.Group("/transactions")
	transactions.Use(s.middleware)
	transactions.POST("", s.addTransaction)
	transactions.GET("", s.transactions)
	transactions.PATCH("/:id", s.updateTransaction)
	transactions.DELETE("/:id", s.deleteTransaction)

	categories := api.Group("/categories")
	categories.Use(s.middleware)
	categories.GET("", s.categories)

	budgets := api.Group("/budgets")
	budgets.Use(s.middleware)
	budgets.POST("", s.addBudget)
	budgets.GET("", s.budgets)

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

func (s *Server) isValidUser(data UserFormat) (bool, string) {
	if !latinRegexp.MatchString(data.Username) || !latinRegexp.MatchString(data.Password) {
		return false, "Имя пользователя и пароль должны состоять только из латинских букв, цифр и особых символов"
	}

	if len(data.Username) < 4 {
		return false, "Имя пользователя должно быть минимум из 4 символов"
	}

	if len(data.Username) > 20 {
		return false, "Имя пользователя должно быть максимум из 20 символов"
	}

	if len(data.Password) < 4 {
		return false, "Пароль должнен быть минимум из 4 символов"
	}

	if len(data.Password) > 30 {
		return false, "Пароль должнен быть максимум из 30 символов"
	}

	return true, ""
}

// @Summary 		Регистрация нового пользователя
// @Description  	Создает аккаунт в системе.
// @Tags         	auth
// @Accept       	json
// @Produce      	json
// @Param        	input body UserFormat true 	"Логин и пароль для регистрации"
// @Success      	201  {object}  MessageResponse 	"Успешная регистрация"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	409  {object}  ErrorResponse 	"Данное имя пользователя уже используется"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка регистрации пользователя"
// @Router       	/auth/register [post]
func (s *Server) register(c *gin.Context) {
	var data UserFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	if isValid, errString := s.isValidUser(data); !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": errString})
		return
	}

	if err := s.authService.Register(data.Username, data.Password); err != nil {
		log.Println(err)
		if errors.Is(err, apperr.NotUnique) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Данное имя пользователя уже используется"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации пользователя"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Успешная регистрация"})
}

// @Summary 		Проверка возможности регистрации пользователя с указанными данными
// @Description  	Проверяет данны по внутренним параметрам, а также уникальности имени пользователя
// @Tags         	auth
// @Accept       	json
// @Produce      	json
// @Param        	input body UserFormat true 	"Логин и пароль для регистрации"
// @Success      	200  {object}  MessageResponse 	"Регистрация возможна"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	400  {object}  ErrorResponse 	"Имя пользователя и пароль должны состоять только из латинских букв, цифр и особых символов"
// @Failure      	400  {object}  ErrorResponse 	"Имя пользователя должно быть минимум из 4 символов"
// @Failure      	400  {object}  ErrorResponse 	"Имя пользователя должно быть максимум из 20 символов"
// @Failure      	400  {object}  ErrorResponse 	"Пароль должнен быть минимум из 4 символов"
// @Failure      	400  {object}  ErrorResponse 	"Пароль должнен быть максимум из 30 символов"
// @Failure      	409  {object}  ErrorResponse 	"Данное имя пользователя уже используется"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при проверке данных"
// @Router       	/auth/register/is-valid [post]
func (s *Server) isValid(c *gin.Context) {
	var data UserFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	isValid, errString := s.isValidUser(data)

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": errString})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Регистрация возможна"})
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
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при создании транзакции"
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

	_, err = s.transactionService.AddTransaction(id, data.Amount, data.Category, data.Description, date)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании транзакции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешно"})
}

// @Summary 		Обновление информации о транзакции
// @Description  	Обновляет информацию о транзакции на новую
// @Tags         	transactions
// @Accept       	json
// @Produce      	json
// @Param			id 		path 		int 					true 	"ID транзакции"
// @Param        	input 	body 		UpdateTransactionFormat true "Новая информация о транзакции"
// @Success      	200  	{object}  	MessageResponse "Успешно"
// @Failure      	400  	{object}  	ErrorResponse 	"Неверно заполнены поля"
// @Failure      	400  	{object}  	ErrorResponse 	"Неверно указано ID транзакции"
// @Failure      	401  	{object}  	ErrorResponse 	"Ошибка авторизации"
// @Failure      	404  	{object}  	ErrorResponse 	"Данная транзакция отсутствует"
// @Failure      	500  	{object}  	ErrorResponse 	"Ошибка при обновлении транзакции"
// @Router       	/transactions/{id} [patch]
func (s *Server) updateTransaction(c *gin.Context) {
	var data UpdateTransactionFormat
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

	userID := c.GetUint("jwt")

	paramID := c.Param("id")
	transactionID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно указано ID транзакции"})
		return
	}

	success, err := s.transactionService.UpdateTransaction(userID, uint(transactionID), data.Amount, data.CategoryID, data.Description, date)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении транзакции"})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данная транзакция отсутствует"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешно"})
}

// @Summary 		Удаление транзакции
// @Description  	Удаление транзакции
// @Tags         	transactions
// @Accept       	json
// @Produce      	json
// @Param			id		path 		int 	true 	"ID транзакции"
// @Success      	200  	{object}  	MessageResponse "Успешно"
// @Failure      	400  	{object}  	ErrorResponse 	"Неверно указано ID транзакции"
// @Failure      	401  	{object}  	ErrorResponse 	"Ошибка авторизации"
// @Failure      	404  	{object}  	ErrorResponse 	"Данная транзакция отсутствует"
// @Failure      	500  	{object}  	ErrorResponse 	"Ошибка при удалениии транзакции"
// @Router       	/transactions/{id} [delete]
func (s *Server) deleteTransaction(c *gin.Context) {
	id := c.GetUint("jwt")

	paramID := c.Param("id")
	transactionID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно указано ID транзакции"})
		return
	}

	success, err := s.transactionService.DeleteTransaction(id, uint(transactionID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении транзакции"})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данная транзакция отсутствует"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешно"})
}

// @Summary 		Получение списка транзакций
// @Description  	Возвращает список транзакций выполненных авторизованным пользователем
// @Tags         	transactions
// @Produce      	json
// @Param        	query query TransactionsFilter true "Параметры транзакций"
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

	var from *time.Time
	if data.From != "" {
		parsedFrom, err := time.Parse("2006-01-02", data.From)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнена дата начала периода"})
			return
		}
		from = &parsedFrom
	}

	var to *time.Time
	if data.To != "" {
		parsedTo, err := time.Parse("2006-01-02", data.To)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнена дата конца периода"})
			return
		}
		to = &parsedTo
	}

	id := c.GetUint("jwt")

	result, err := s.transactionService.Transactions(id, data.Limit, data.Page, from, to, data.Type, data.Categories)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении транзакций"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// @Summary 		Получение списка системных категорий
// @Description  	Возвращает список системных категорий
// @Tags         	categories
// @Produce      	json
// @Success      	200  {object}  CategoriesResponse 		"Успешное получение системных категорий"
// @Failure      	401  {object}  ErrorResponse 			"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 			"Ошибка при получении категории"
// @Router       	/system-categories [get]
// @Ignore
func (s *Server) systemCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"result": s.categoryService.SystemCategories()})
}

// @Summary 		Получение списка пользовательских категорий
// @Description  	Возвращает список пользовательских категорий
// @Tags         	categories
// @Produce      	json
// @Success      	200  {object}  CategoriesResponse 		"Успешное получение пользовательских категорий"
// @Failure      	401  {object}  ErrorResponse 			"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 			"Ошибка при получении категории"
// @Router       	/user-categories [get]
// @Ignore
func (s *Server) userCategories(c *gin.Context) {
	id := c.GetUint("jwt")

	result, err := s.categoryService.UserCategories(id)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении категорий"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// @Summary 		Получение списка категорий
// @Description  	Возвращает список пользовательских и системных категорий
// @Tags         	categories
// @Produce      	json
// @Success      	200  {object}  CategoriesResponse 		"Успешное получение категорий"
// @Failure      	401  {object}  ErrorResponse 			"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 			"Ошибка при получении категории"
// @Router       	/categories [get]
func (s *Server) categories(c *gin.Context) {
	id := c.GetUint("jwt")

	result, err := s.categoryService.Categories(id)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении категорий"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// @Summary 		Создание бюджета
// @Description  	Создает бюджет от имени авторизированного пользователя
// @Tags         	budgets
// @Accept       	json
// @Produce      	json
// @Param        	input body AddBudgetFormat true "Категория и лимит"
// @Success      	200  {object}  MessageResponse 	"Успешно"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	401  {object}  ErrorResponse 	"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при создании бюджета"
// @Router       	/budgets [post]
func (s *Server) addBudget(c *gin.Context) {
	var data AddBudgetFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	id := c.GetUint("jwt")

	if err := s.budgetService.AddBudget(id, data.CategoryID, data.Limit); err != nil {
		log.Println(err)
		if errors.Is(err, apperr.NotUnique) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Бюджет с данной категорией уже существует"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании бюджета"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешно"})
}

// @Summary 		Получение информации по бюджетам
// @Description  	Возвращает список бюджет и текущие затраты по ним
// @Tags         	budgets
// @Accept       	json
// @Produce      	json
// @Param        	query query BudgetsFilter true "Дата в котором рассамтривается бюджет, из нее берется год и месяц"
// @Success      	200  {object}  BudgetsResponse 			"Успешное получение бюджетов"
// @Failure      	400  {object}  ErrorResponse 			"Неверно заполнена дата"
// @Failure      	401  {object}  ErrorResponse 			"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 			"Ошибка при получении бюджетов"
// @Router       	/budgets [get]
func (s *Server) budgets(c *gin.Context) {
	var data BudgetsFilter
	if err := c.ShouldBindQuery(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнена дата"})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнена дата"})
		return
	}

	id := c.GetUint("jwt")

	result, err := s.budgetService.Budgets(id, parsedDate)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении бюджетов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
