package main

import (
	"log"
	"mobilepay"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	port := "8080"

	gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	r := gin.Default()

	var LogLevel mobilepay.LeveledLoggerInterface = &mobilepay.LeveledLogger{
		Level: mobilepay.LevelDebug,
	}

	config := &mobilepay.Config{
		HTTPClient: http.DefaultClient,
		Logger:     LogLevel,
	}

	client := mobilepay.NewClient(
		os.Getenv("MOBILEPAY_CLIENT_ID"),
		os.Getenv("MOBILEPAY_API_KEY"),
		config,
	)

	PaymentsResource{mp: client}.Routes(r)

	r.Run(":" + port)
}
