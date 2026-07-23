package server

import (
	"errors"
	"finstat/internal/apperr"
	"finstat/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryFormat struct {
	Name string `json:"name" binding:"required"`
}

type CategoriesResponse struct {
	Result []service.Category `json:"result"`
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

	id, err := s.services.Category.InsertCategory(userID, data.Name)
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, apperr.NotUnique):
			c.JSON(http.StatusConflict, gin.H{"error": "Категория с таким именем уже есть"})
		default:
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

	success, err := s.services.Category.UpdateCategory(userID, uint(categoryID), data.Name)

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

	success, err := s.services.Category.DeleteCategory(userID, uint(budgetID))
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

	result, err := s.services.Category.Categories(userID)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении категорий"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
