package server

import (
	"finstat/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type InsertTransactionFormat struct {
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

type TransactionsResponse struct {
	Result []service.Transaction `json:"result"`
}

// @Summary 		Создание транзакции
// @Description  	Создает транзакцию от имени авторизированного пользователя
// @Tags         	transactions
// @Accept       	json
// @Produce      	json
// @Param        	input body InsertTransactionFormat true "Информация о транзакции"
// @Success      	200								"Успешное создание транзакции"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	401  {object}  ErrorResponse 	"Ошибка авторизации"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при создании транзакции"
// @Router       	/transactions [post]
// @Security     	Auth
func (s *Server) addTransaction(c *gin.Context) {
	var data InsertTransactionFormat
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

	_, err = s.services.Transaction.InsertTransaction(userID, data.Amount, data.CategoryID, data.Description, date)
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

	success, err := s.services.Transaction.UpdateTransaction(userID, uint(transactionID), data.Amount, data.CategoryID, data.Description, date)
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

	success, err := s.services.Transaction.DeleteTransaction(userID, uint(transactionID))
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

	result, err := s.services.Transaction.Transactions(userID, data.Limit, data.Page, from, to, data.Type, data.Categories)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении транзакций"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
