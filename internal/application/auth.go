package application

import (
	"database/sql"
	"fmt"
	"github.com/btnbrd/avitoshop/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecretkey")

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

//func auth(c *gin.Context) {
//	c.String(http.StatusOK, "auth")
//}

func (app *APIServer) AuthHandler(c *gin.Context) {
	log.Println("authhandler")
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		log.Println("Invalid request:", err) // Отладочная информация
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Неверный запрос"})
		return
	}

	user, err := app.getUserByUsername(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			user, err = app.createUser(req.Username, req.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка создания пользователя"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
			return
		}
	} else {
		// Сравниваем пароль
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
			log.Println("Invalid credentials") // Отладочная информация
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Неверные учетные данные"})
			return
		}
	}

	fmt.Printf("%+v\n", user)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка генерации токена"})
		return
	}
	c.JSON(http.StatusOK, models.AuthResponse{Token: tokenString})
}
func (app *APIServer) getUserByUsername(username string) (*models.User, error) {
	row := app.db.QueryRow("SELECT id, username, password, coins FROM users WHERE username=$1", username)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Coins)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (app *APIServer) createUser(username string, password string) (*models.User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var id int
	err = app.db.QueryRow(
		"INSERT INTO users(username, password, coins) VALUES($1, $2, 1000) RETURNING id",
		username, string(hashedPass)).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:       id,
		Username: username,
		Password: string(hashedPass),
		Coins:    1000,
	}, nil
}

func (app *APIServer) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "No Authorization header provided"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Invalid Authorization header format"})
			c.Abort()
			return
		}
		tokenStr := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
