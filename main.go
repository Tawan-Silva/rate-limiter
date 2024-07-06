package main

import (
	"log"
	"os"
	"ratelimiter/configs"
	_ "ratelimiter/docs"
	"ratelimiter/middleware"
	"ratelimiter/pkg/ratelimiter"
	"ratelimiter/server"
)

// @title           Rate Limiter API Example
// @version         1.0
// @description     Rate Limiter API with Redis
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   Tawan Silva
// @contact.url    https://www.linkedin.com/in/tawan-silva-684b581b7/
// @contact.email tawan.tls43@gmail.com
//
// @license.name   Rate Limiter License
// @license.url   http://www.ratelimiter.com.br
//
// @host      localhost:8080
// @BasePath  /
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name API_KEY
// @type apiKey
func main() {
	_, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	store := ConnectToRedis()
	if err := store.Ping(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")

	rateLimiter := ratelimiter.NewRateLimiter(store)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiter)

	s := server.NewServer(rateLimiterMiddleware)
	log.Println("Starting the server...")
	s.Start()
	log.Println("Server started successfully")
}

func ConnectToRedis() *ratelimiter.RedisStore {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		redisAddress = "localhost:6379"
	}
	return ratelimiter.NewRedisStore(redisAddress)
}
