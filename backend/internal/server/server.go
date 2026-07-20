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
	JWT_COOKIE_NAME = "refresh_jwt"
	USER_ID_KEY     = "user_id"
	REFRESH_PATH    = "/api/auth"
)

var (
	latinRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// @title Finstat API
// @version 0.9
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
	host               string
	authService        *service.AuthService
	transactionService *service.TransactionService
	categoryService    *service.CategoryService
	budgetService      *service.BudgetService
}

type UserFormat struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type IsUserValidFormat struct {
	Username string `json:"username" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty"`
}

type AddTransactionFormat struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	CategoryID  uint            `json:"category_id" binding:"required,numeric"`
	Description string          `json:"description" binding:"omitempty"`
	Date        string          `json:"date" binding:"required"`
}

type UpdateTransactionFormat struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	CategoryID  uint            `json:"category_id" binding:"required,numeric"`
	Description string          `json:"description" binding:"omitempty"`
	Date        string          `json:"date" binding:"required"`
}

type DeleteTransactionFormat struct {
	ID uint `json:"id" binding:"required, numeric"`
}

type TransactionsFilter struct {
	Limit      uint           `form:"limit" binding:"required,numeric,min=1"`
	Page       uint           `form:"page" binding:"required,numeric,min=1"`
	From       string         `form:"from" binding:"omitempty,datetime=2006-01-02"`
	To         string         `form:"to" binding:"omitempty,datetime=2006-01-02"`
	Type       int            `form:"type" binding:"omitempty,numeric,min=-1,max=1"`
	Categories CategoriesList `form:"categories,parser=encoding.TextUnmarshaler" binding:"omitempty"`
}

type CategoryFormat struct {
	Name string `json:"name" binding:"required"`
}

type AddBudgetFormat struct {
	CategoryID uint            `json:"category_id" binding:"required"`
	Limit      decimal.Decimal `json:"limit" binding:"required"`
}

type UpdateBudgetFormat struct {
	Limit decimal.Decimal `json:"limit" binding:"required"`
}

type BudgetsFilter struct {
	Date string `form:"date" binding:"required,datetime=2006-01-02"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type StringResponse struct {
	Result string `json:"result"`
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

type BudgetResponse struct {
	Result service.Budget `json:"result"`
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

	id, err := s.authService.ID(token)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}

	c.Set(USER_ID_KEY, id)

	c.Next()
}

func (s *Server) isValidUser(data UserFormat) (bool, string) {
	if data.Username == "" {
		return false, "Имя пользователя не может быть пустым"
	}

	if data.Password == "" {
		return false, "Пароль не может быть пустым"
	}

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

// @Summary 		Проверка возможности регистрации пользователя с указанными данными
// @Description  	Проверяет данные по внутренним параметрам, а также уникальности имени пользователя
// @Tags         	auth
// @Accept       	json
// @Produce      	json
// @Param        	input body IsUserValidFormat true 	"Логин и пароль для регистрации"
// @Success      	200								"Регистрация возможна"
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
	var data IsUserValidFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	user := UserFormat{
		Username: data.Username,
		Password: data.Password,
	}

	isValid, errString := s.isValidUser(user)

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": errString})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary 		Регистрация нового пользователя
// @Description  	Создает аккаунт в системе.
// @Tags         	auth
// @Accept       	json
// @Produce      	json
// @Param        	input body UserFormat true 	"Логин и пароль для регистрации"
// @Success      	201 							"Успешная регистрация"
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

	c.JSON(http.StatusCreated, nil)
}

// @Summary 		Авторизация пользователя
// @Description  	Авторизирует пользователя системы и сохраняет jwt-токена в cookies
// @Tags         	auth
// @Accept       	json
// @Produce      	json
// @Param        	input body UserFormat true 		"Логин и пароль для авторизации"
// @Success      	200								"Успешная авторизация"
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

	access, refresh, err := s.authService.Login(data.Username, data.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка авторизации пользователя"})
		return
	}

	c.SetCookie(
		JWT_COOKIE_NAME,
		refresh,
		service.REFRESH_TOKEN_LIFE_TIME,
		REFRESH_PATH,
		s.host,
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"result": access})
}

// @Summary 		Обновление токенов
// @Description  	Обновляет access и refresh токены
// @Tags         	auth
// @Accept       	json
// @Produce      	json
// @Success      	200								"Успешное обновление токенов"
// @Failure      	401	{object}	ErrorResponse	"Пользователь не авторизован"
// @Failure      	500	{object}	ErrorResponse 	"Ошибка обновления токенов"
// @Router       	/auth/refresh [post]
func (s *Server) refresh(c *gin.Context) {
	cookie, err := c.Cookie(JWT_COOKIE_NAME)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	access, refresh, err := s.authService.Refresh(cookie)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления токенов"})
		return
	}

	c.SetCookie(
		JWT_COOKIE_NAME,
		refresh,
		service.REFRESH_TOKEN_LIFE_TIME*24*60*60,
		REFRESH_PATH,
		s.host,
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"result": access})
}

// @Summary 		Деавторизация пользователя
// @Description  	Деавторизирует пользователя системы и удаляет jwt-токены
// @Tags         	auth
// @Produce      	json
// @Success      	200			"Успешная деавторизация"
// @Success			202			"Успешное удаление токена из куки"
// @Failure      	401			"Пользователь не авторизован"
// @Router       	/auth/logout [post]
func (s *Server) logout(c *gin.Context) {
	cookie, err := c.Cookie(JWT_COOKIE_NAME)

	c.SetCookie(
		JWT_COOKIE_NAME,
		"",
		-1,
		REFRESH_PATH,
		s.host,
		true,
		true,
	)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	err = s.authService.Logout(cookie)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusAccepted, gin.H{"error": "Успешное удаление токена из куки"})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary 		Создание транзакции
// @Description  	Создает транзакцию от имени авторизированного пользователя
// @Tags         	transactions
// @Accept       	json
// @Produce      	json
// @Param        	input body AddTransactionFormat true "Информация о транзакции"
// @Success      	200								"Успешное создание транзакции"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	401  {object}  ErrorResponse 	"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при создании транзакции"
// @Router       	/transactions [post]
// @Security     	Auth
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

	userID := c.GetUint(USER_ID_KEY)

	_, err = s.transactionService.AddTransaction(userID, data.Amount, data.CategoryID, data.Description, date)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании транзакции"})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary 		Обновление информации о транзакции
// @Description  	Обновляет информацию о транзакции на новую
// @Tags         	transactions
// @Accept       	json
// @Produce      	json
// @Param			id 		path 		int 					true 	"ID транзакции"
// @Param        	input 	body 		UpdateTransactionFormat true 	"Новая информация о транзакции"
// @Success      	200									"Успешное обновление транзакции"
// @Failure      	400  	{object}  	ErrorResponse 	"Неверно заполнены поля"
// @Failure      	400  	{object}  	ErrorResponse 	"Неверно указано ID транзакции"
// @Failure      	401  	{object}  	ErrorResponse 	"Ошибка авторизации"
// @Failure      	404  	{object}  	ErrorResponse 	"Данная транзакция отсутствует"
// @Failure      	500  	{object}  	ErrorResponse 	"Ошибка при обновлении транзакции"
// @Router       	/transactions/{id} [patch]
// @Security     	Auth
func (s *Server) updateTransaction(c *gin.Context) {
	paramID := c.Param("id")
	transactionID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно указано ID транзакции"})
		return
	}

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

	userID := c.GetUint(USER_ID_KEY)

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

	c.JSON(http.StatusOK, nil)
}

// @Summary 		Удаление транзакции
// @Description  	Удаляет транзакцию авторизованного пользователя
// @Tags         	transactions
// @Accept       	json
// @Produce      	json
// @Param			id		path 		int 	true 	"ID транзакции"
// @Success      	200									"Успешное удаление транзакции"
// @Failure      	400  	{object}  	ErrorResponse 	"Неверно указано ID транзакции"
// @Failure      	401  	{object}  	ErrorResponse 	"Ошибка авторизации"
// @Failure      	404  	{object}  	ErrorResponse 	"Данная транзакция отсутствует"
// @Failure      	500  	{object}  	ErrorResponse 	"Ошибка при удалениии транзакции"
// @Router       	/transactions/{id} [delete]
// @Security     	Auth
func (s *Server) deleteTransaction(c *gin.Context) {
	userID := c.GetUint(USER_ID_KEY)

	paramID := c.Param("id")
	transactionID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Неверно указано ID транзакции"})
		return
	}

	success, err := s.transactionService.DeleteTransaction(userID, uint(transactionID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении транзакции"})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данная транзакция отсутствует"})
		return
	}

	c.JSON(http.StatusOK, nil)
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
// @Security     	Auth
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

	userID := c.GetUint(USER_ID_KEY)

	result, err := s.transactionService.Transactions(userID, data.Limit, data.Page, from, to, data.Type, data.Categories)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении транзакций"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// @Summary 		Создание пользовательской категории
// @Description  	Создает пользовательскую категорию и возвращает ее ID
// @Tags         	categories
// @Accept       	json
// @Produce      	json
// @Param        	input body CategoryFormat true "Имя категории (Больше 2 символов)"
// @Success      	200								"Успешное создание пользовательской категории"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнено имя категории"
// @Failure      	401  {object}  ErrorResponse 	"Ошибка авторизации"
// @Failure      	409  {object}  ErrorResponse 	"Категория с таким именем уже есть"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при создании категории"
// @Router       	/categories [post]
// @Security     	Auth
func (s *Server) addCategory(c *gin.Context) {
	var data CategoryFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнено имя категории"})
		return
	}

	if len(data.Name) <= 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя категории должно быть больше 2 символов"})
		return
	}

	userID := c.GetUint(USER_ID_KEY)

	id, err := s.categoryService.AddCategory(userID, data.Name)
	if err != nil {
		log.Println(err)
		if errors.Is(err, apperr.NotUnique) {
			c.JSON(http.StatusConflict, gin.H{"error": "Категория с таким именем уже есть"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании категории"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary 		Обновление категории
// @Description  	Обновляет название категории авторизированного пользователя
// @Tags         	categories
// @Accept       	json
// @Produce      	json
// @Param			id		path 		int 				true	"ID категории"
// @Param        	input 	body 		CategoryFormat 		true	"Новое имя категории"
// @Success      	200 							"Успешное обновление категории"
// @Failure      	400  {object}  ErrorResponse 	"Неверно указано ID бюджета"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнен новый лимит"
// @Failure      	401  {object}  ErrorResponse 	"Ошибка авторизации"
// @Failure      	404  {object}  ErrorResponse 	"Данный бюджет отсутствует"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при создании бюджета"
// @Router       	/categories/{id} [patch]
// @Security     	Auth
func (s *Server) updateCategory(c *gin.Context) {
	paramID := c.Param("id")
	categoryID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно указано ID категории"})
		return
	}

	var data CategoryFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнено новое имя категории"})
		return
	}

	userID := c.GetUint(USER_ID_KEY)

	success, err := s.categoryService.UpdateCategory(userID, uint(categoryID), data.Name)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении категории"})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данная категория отсутствует"})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary 		Удаление категории
// @Description  	Удаляет категорию авторизованного пользователя
// @Tags         	categories
// @Produce      	json
// @Param			id		path 		int 	true 	"ID категории"
// @Success      	200									"Успешное удаление категории"
// @Failure      	400  	{object}  	ErrorResponse 	"Неверно указано ID категории"
// @Failure      	401  	{object}  	ErrorResponse 	"Ошибка авторизации"
// @Failure      	404  	{object}  	ErrorResponse 	"Данная категория отсутствует"
// @Failure      	500  	{object}  	ErrorResponse 	"Ошибка при удалениии категории"
// @Router       	/categories/{id} [delete]
// @Security     	Auth
func (s *Server) deleteCategory(c *gin.Context) {
	paramID := c.Param("id")
	budgetID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Неверно указано ID категории"})
		return
	}

	userID := c.GetUint(USER_ID_KEY)

	success, err := s.categoryService.DeleteCategory(userID, uint(budgetID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении категории"})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данная категория отсутствует"})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary 		Получение списка категорий
// @Description  	Возвращает список пользовательских и системных категорий
// @Tags         	categories
// @Produce      	json
// @Success      	200  {object}  CategoriesResponse 		"Успешное получение категорий"
// @Failure      	401  {object}  ErrorResponse 			"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 			"Ошибка при получении категории"
// @Router       	/categories [get]
// @Security     	Auth
func (s *Server) categories(c *gin.Context) {
	userID := c.GetUint(USER_ID_KEY)

	result, err := s.categoryService.Categories(userID)

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
// @Success      	200								"Успешное создание бюджета"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	401  {object}  ErrorResponse 	"Ошибка авторизации"
// @Failure      	409  {object}  ErrorResponse 	"Бюджет с данной категорией уже существует"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при создании бюджета"
// @Router       	/budgets [post]
// @Security     	Auth
func (s *Server) addBudget(c *gin.Context) {
	var data AddBudgetFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	userID := c.GetUint(USER_ID_KEY)

	if err := s.budgetService.AddBudget(userID, data.CategoryID, data.Limit); err != nil {
		log.Println(err)
		if errors.Is(err, apperr.NotUnique) {
			c.JSON(http.StatusConflict, gin.H{"error": "Бюджет с данной категорией уже существует"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании бюджета"})
		}
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary 		Обновление бюджета
// @Description  	Обновляет лимит бюджета авторизированного пользователя
// @Tags         	budgets
// @Accept       	json
// @Produce      	json
// @Param			id		path 		int 				true	"ID бюджета"
// @Param        	input 	body 		UpdateBudgetFormat 	true	"Новый лимит"
// @Success      	200 							"Успешное обновление бюджета"
// @Failure      	400  {object}  ErrorResponse 	"Неверно указано ID бюджета"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнен новый лимит"
// @Failure      	401  {object}  ErrorResponse 	"Ошибка авторизации"
// @Failure      	404  {object}  ErrorResponse 	"Данный бюджет отсутствует"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при создании бюджета"
// @Router       	/budgets/{id} [patch]
// @Security     	Auth
func (s *Server) updateBudget(c *gin.Context) {
	paramID := c.Param("id")
	budgetID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно указано ID бюджета"})
		return
	}

	var data UpdateBudgetFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнен новый лимит"})
		return
	}

	userID := c.GetUint(USER_ID_KEY)

	success, err := s.budgetService.UpdateBudget(userID, uint(budgetID), data.Limit)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении бюджета"})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данный бюджет отсутствует"})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary 		Удаление бюджета
// @Description  	Удаляет бюджет авторизованного пользователя
// @Tags         	budgets
// @Produce      	json
// @Param			id		path 		int 	true 	"ID бюджета"
// @Success      	200									"Успешное удаление бюджета"
// @Failure      	400  	{object}  	ErrorResponse 	"Неверно указано ID бюджета"
// @Failure      	401  	{object}  	ErrorResponse 	"Ошибка авторизации"
// @Failure      	404  	{object}  	ErrorResponse 	"Данный бюджет отсутствует"
// @Failure      	500  	{object}  	ErrorResponse 	"Ошибка при удалениии бюджета"
// @Router       	/budgets/{id} [delete]
// @Security     	Auth
func (s *Server) deleteBudget(c *gin.Context) {
	paramID := c.Param("id")
	budgetID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Неверно указано ID бюджета"})
		return
	}

	userID := c.GetUint(USER_ID_KEY)

	success, err := s.budgetService.DeleteBudget(userID, uint(budgetID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении бюджета"})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Данный бюджет отсутствует"})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary 		Получение информации по бюджетам
// @Description  	Возвращает список бюджет и текущие затраты по ним
// @Tags         	budgets
// @Accept       	json
// @Produce      	json
// @Param        	query query BudgetsFilter true "Дата в котором рассамтриваются бюджеты, из нее берется год и месяц"
// @Success      	200  {object}  BudgetsResponse 			"Успешное получение бюджетов"
// @Failure      	400  {object}  ErrorResponse 			"Неверно заполнена дата"
// @Failure      	401  {object}  ErrorResponse 			"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 			"Ошибка при получении бюджетов"
// @Router       	/budgets [get]
// @Security     	Auth
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

	userID := c.GetUint(USER_ID_KEY)

	result, err := s.budgetService.Budgets(userID, parsedDate)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении бюджетов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// @Summary 		Получение информации о бюджете на данную категорию
// @Description  	Возвращает бюджет и текущие затраты
// @Tags         	budgets
// @Accept       	json
// @Produce      	json
// @Param			id		path 		int 				true	"ID категории"
// @Param        	query 	query 		BudgetsFilter 		true 	"Дата в котором рассамтривается бюджет, из нее берется год и месяц"
// @Success      	200  {object}  BudgetResponse 			"Успешное получение бюджетов"
// @Failure      	400  {object}  ErrorResponse 			"Неверно заполнена дата"
// @Failure      	401  {object}  ErrorResponse 			"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 			"Ошибка при получении бюджетов"
// @Router       	/budgets/category/{id} [get]
// @Security     	Auth
func (s *Server) budgetByCategory(c *gin.Context) {
	paramID := c.Param("id")
	categoryID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно указано ID категории"})
		return
	}

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

	userID := c.GetUint(USER_ID_KEY)

	result, err := s.budgetService.BudgetByCategory(userID, uint(categoryID), parsedDate)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении бюджетов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
