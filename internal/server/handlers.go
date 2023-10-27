package server

import (
	"crypto/hmac"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/sha3"
	"melvadTestTask/internal/interfaces"
	"net/http"
)

type RedisRequest struct {
	Key   string `json:"key" binding:"required"`
	Value int    `json:"value" binding:"required"`
}

type SignRequest struct {
	Text string `json:"text" binding:"required"`
	Key  string `json:"key" binding:"required"`
}

type UserRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required"`
}

type UserResponse struct {
	ID int `json:"id"`
}

func handleRedisIncrement(redis interfaces.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RedisRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Инкрементировать значение
		newVal, err := redis.IncrBy(req.Key, int64(req.Value)).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"value": newVal})
	}
}

func handleHMACSHA512() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		key := []byte(req.Key)
		text := []byte(req.Text)
		h := hmac.New(sha3.New512, key)
		h.Write(text)
		signature := hex.EncodeToString(h.Sum(nil))

		c.JSON(http.StatusOK, gin.H{"signature": signature})
	}
}

func handlePostgresUser(db interfaces.Postgres) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Создать таблицу users, если она не существует
		createTableSQL := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name TEXT,
            age INT
        )
    `
		_, err := db.Exec(createTableSQL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Добавить пользователя в базу данных  и получить его ключ
		insertUserSQL := `INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id`
		var userID int
		err = db.Get(&userID, insertUserSQL, req.Name, req.Age)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, UserResponse{ID: userID})
	}
}
