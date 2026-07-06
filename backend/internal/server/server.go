package server

import (
	"log"
	"net/http"

	ewrap "auth.my-financials/internal/lib"
	"auth.my-financials/internal/repository"
	"auth.my-financials/internal/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	postgresClient *repository.Client
	jwt_secret     []byte
}

type UserData struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}

func NewServer(postgres_creds repository.Credentials, jwt_secret string) (*Server, error) {
	db_client, err := repository.NewClient(postgres_creds)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't create server", err)
	}

	return &Server{
		postgresClient: db_client,
		jwt_secret:     []byte(jwt_secret),
	}, nil
}

func (s *Server) Start(port string) {
	router := gin.Default()

	authGroup := router.Group("/auth")
	authGroup.POST("/register", s.register)
	authGroup.POST("/login", s.login)

	router.Run(":" + port)
}

func (s *Server) register(c *gin.Context) {
	var data UserData
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации хэша"})
		return
	}

	_, err = s.postgresClient.InsertUser(data.Username, string(hashedPassword))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Успешная регистрация"})
}

func (s *Server) login(c *gin.Context) {
	var data UserData
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверно заполнены поля"})
		return
	}

	user, err := s.postgresClient.User(data.Username)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(data.Password)); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	newToken, err := token.NewToken(user.ID, s.jwt_secret)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка авторизации"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешная авторизация", "token": newToken})
}
