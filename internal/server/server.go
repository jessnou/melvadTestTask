package server

import (
	"github.com/gin-gonic/gin"
	"melvadTestTask/internal/config"
)

func StartServer(app *config.App) {
	router := setupRouter(app)
	router.Run(":8080")
}
func setupRouter(app *config.App) *gin.Engine {
	router := gin.Default()
	router.POST("/redis/incr", handleRedisIncrement(app.RedisClient))
	router.POST("/sign/hmacsha512", handleHMACSHA512())
	router.POST("/postgres/users", handlePostgresUser(app.DB))
	return router
}
