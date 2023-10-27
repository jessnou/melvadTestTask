package db

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"log"
	"melvadTestTask/internal/config"
	"os"

	_ "github.com/lib/pq"
)

func GetDBConn() (*sqlx.DB, *redis.Client, error) {
	redisHost := os.Args[1]
	redisPort := os.Args[2]
	conf, err := config.LoadConfig(".")
	if err != nil {
		if err != nil {
			log.Fatal("cannot load config", err)
		}
	}

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: fmt.Sprintf("%s:%d", conf.RedisHost, conf.RedisPort),
	//})
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", conf.DbUser, conf.DbPassword, conf.DbName))
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL:", err)
	}

	return db, redisClient, nil
}
