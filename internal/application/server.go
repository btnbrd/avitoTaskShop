package application

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"math/rand/v2"
	"net/http"
)

type APIServer struct {
	r  *gin.Engine
	db *sql.DB
}

func NewServer(db *sql.DB) *APIServer {
	r := gin.Default()
	return &APIServer{r: r, db: db}
}

func (a *APIServer) Run() error {
	a.r.GET("/g/:id", PathParameter)
	a.r.GET("/", Simple)
	//a.r.GET("/info", getInfo)
	//a.r.POST("/sendCoin", sendCoin)
	//a.r.GET("/buy/:item", buyItem)
	a.r.POST("/auth", a.authHandler)
	authGroup := a.r.Group("/")
	authGroup.Use(a.AuthMiddleware())
	{
		authGroup.GET("/info", a.infoHandler)
		authGroup.POST("/sendCoin", a.sendCoinHandler)
		authGroup.GET("/buy/:item", a.buyHandler)
	}

	return a.r.Run(":8080")
}

// Функция-обработчик теперь принимает *gin.Context
func Simple(c *gin.Context) {
	c.String(http.StatusOK, "hello %d", rand.Uint())
}

func PathParameter(c *gin.Context) {
	id := c.Param("id") // Получаем параметр из пути
	c.JSON(http.StatusOK, gin.H{"user_id": id})
	//c.String()
}
