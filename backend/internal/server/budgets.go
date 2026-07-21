package server

import (
	"errors"
	"finstat/internal/apperr"
	"finstat/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type InsertBudgetFormat struct {
	CategoryID uint            `json:"category_id" binding:"required"`
	Limit      decimal.Decimal `json:"limit" binding:"required"`
}

type UpdateBudgetFormat struct {
	Limit decimal.Decimal `json:"limit" binding:"required"`
}

type BudgetsFilter struct {
	Date string `form:"date" binding:"required,datetime=2006-01-02"`
}

type BudgetsResponse struct {
	Result []service.Budget `json:"result"`
}

type BudgetResponse struct {
	Result service.Budget `json:"result"`
}

// @Summary 		Создание бюджета
// @Description  	Создает бюджет от имени авторизированного пользователя
// @Tags         	budgets
// @Accept       	json
// @Produce      	json
// @Param        	input body InsertBudgetFormat true "Категория и лимит"
// @Success      	200								"Успешное создание бюджета"
// @Failure      	400  {object}  ErrorResponse 	"Неверно заполнены поля"
// @Failure      	401  {object}  ErrorResponse 	"Ошибка авторизации"
// @Failure      	409  {object}  ErrorResponse 	"Бюджет с данной категорией уже существует"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка при создании бюджета"
// @Router       	/budgets [post]
// @Security     	Auth
func (s *Server) addBudget(c *gin.Context) {
	var data InsertBudgetFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	userID := c.GetUint(USER_ID_KEY)

	if err := s.services.Budget.InsertBudget(userID, data.CategoryID, data.Limit); err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, apperr.NotUnique):
			c.JSON(http.StatusConflict, gin.H{"error": "Бюджет с данной категорией уже существует"})
		default:
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

	success, err := s.services.Budget.UpdateBudget(userID, uint(budgetID), data.Limit)

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

	success, err := s.services.Budget.DeleteBudget(userID, uint(budgetID))
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
// @Failure      	404  {object}  ErrorResponse 			"Бюджеты отсутствуют"
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

	result, err := s.services.Budget.Budgets(userID, parsedDate)

	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, apperr.NoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": "Бюджеты отсутствуют"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении бюджетов"})
		}
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
// @Failure      	404  {object}  ErrorResponse 			"Бюджет с данной категорией отсутствует"
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

	result, err := s.services.Budget.BudgetByCategory(userID, uint(categoryID), parsedDate)

	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, apperr.NoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": "Бюджет с данной категорией отсутствует"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении бюджетов"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
