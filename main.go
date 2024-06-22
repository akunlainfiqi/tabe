package main

import (
	"saas-billing/config"
	"saas-billing/presentation/httpgin"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	router := httpgin.New()
	router.Run(":" + config.PORT)
}
