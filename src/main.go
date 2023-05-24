package main

import (
	"fmt"
	"os"

	"ta/backend/src/router"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	start()
}

func start() {
	appRouter := router.SetupRouter()
	router.LoadRouter(appRouter)
	serverPort := os.Getenv("SERVER_PORT")
	appRouter.Run(fmt.Sprintf(":%s", serverPort))
}
