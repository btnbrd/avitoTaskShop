package application

import (
	"github.com/btnbrd/avitoshop/internal/storage"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	r  *gin.Engine
	db storage.DBInterface
}

func NewServer(db storage.DBInterface) *APIServer {
	r := gin.Default()
	return &APIServer{r: r, db: db}
}

func (a *APIServer) Run() error {

	a.r.POST("/auth", a.AuthHandler)
	authGroup := a.r.Group("/")
	authGroup.Use(a.AuthMiddleware())
	{
		authGroup.GET("/info", a.infoHandler)
		authGroup.POST("/sendCoin", a.sendCoinHandler)
		authGroup.GET("/buy/:item", a.buyHandler)
	}

	return a.r.Run(":8080")
}
