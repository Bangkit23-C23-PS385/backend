package router

import (
	"backend/docs"
	v1 "backend/src/controller/v1"
	"backend/src/database"
	"backend/src/util/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var master *gorm.DB
var slave *gorm.DB

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	return router
}

func LoadRouter(router *gin.Engine) {
	initSwagger(router)
	initMaster()

	// add router
	v1.RouteLoader(router, database.DB{
		Master: master,
	})
}

func initMaster() {
	master = database.DBMaster()
}

func initSwagger(router *gin.Engine) {
	_ = godotenv.Load()

	docs.SwaggerInfo.Title = "MediCare API"
	docs.SwaggerInfo.Description = "MediCare API Documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("API_ORIGIN_URL")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	if os.Getenv("ENV") != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
