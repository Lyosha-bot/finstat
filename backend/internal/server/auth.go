package server

import (
	"errors"
	"finstat/internal/apperr"
	"finstat/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserFormat struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type IsUserValidFormat struct {
	Username string `json:"username" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty"`
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

	if err := s.services.Auth.Register(data.Username, data.Password); err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, apperr.NotUnique):
			c.JSON(http.StatusConflict, gin.H{"error": "Данное имя пользователя уже используется"})
		default:
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
// @Failure      	400  {object}  ErrorResponse 	"Нверное имя пользователя или пароль"
// @Failure      	500  {object}  ErrorResponse 	"Ошибка авторизации пользователя"
// @Router       	/auth/login [post]
func (s *Server) login(c *gin.Context) {
	var data UserFormat
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	access, refresh, err := s.services.Auth.Login(data.Username, data.Password)
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, apperr.NoRows):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Нверное имя пользователя или пароль"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка авторизации пользователя"})
		}
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

	access, refresh, err := s.services.Auth.Refresh(cookie)
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

	err = s.services.Auth.Logout(cookie)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusAccepted, gin.H{"error": "Успешное удаление токена из куки"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
