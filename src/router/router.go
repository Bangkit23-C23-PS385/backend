package router

import (
	"os"
	"sync"
	"ta/backend/docs"
	v1 "ta/backend/src/controller/v1"
	"ta/backend/src/database"
	"ta/backend/src/util/middleware"

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

	wg := &sync.WaitGroup{}

	wg.Add(2)
	go initMaster(wg)
	go initSlave(wg)
	wg.Wait()

	// add router
	v1.RouteLoader(router, database.DB{
		Master: master,
		Slave:  slave,
	})
}

func initMaster(wg *sync.WaitGroup) {
	defer wg.Done()
	master = database.DBMaster()
}

func initSlave(wg *sync.WaitGroup) {
	defer wg.Done()
	slave = database.DBSlave()
}

func initSwagger(router *gin.Engine) {
	_ = godotenv.Load()

	docs.SwaggerInfo.Title = "AIVue API"
	docs.SwaggerInfo.Description = "AIVue API Documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("API_ORIGIN_URL")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	if os.Getenv("ENV") != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
