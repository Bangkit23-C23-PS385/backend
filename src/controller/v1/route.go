package v1

import (
	authCtrl "backend/src/controller/v1/auth"
	predictCtrl "backend/src/controller/v1/predict"
	"backend/src/database"
	authRepo "backend/src/repository/v1/auth"
	verifRepo "backend/src/repository/v1/verification"
	authSvc "backend/src/service/v1/auth"
	predictSvc "backend/src/service/v1/predict"
	verifSvc "backend/src/service/v1/verification"

	"github.com/gin-gonic/gin"
)

func RouteLoader(router *gin.Engine, db database.DB) {
	verifRepository := verifRepo.NewRepository(db)
	authRepository := authRepo.NewRepository(db)

	predictService := predictSvc.NewService()
	verifService := verifSvc.NewService(verifRepository)
	authService := authSvc.NewService(verifService, authRepository)

	predictController := predictCtrl.NewController(predictService)
	authController := authCtrl.NewController(authService)

	v1 := router.Group("v1")

	auth(v1, authController)
	predict(v1, predictController)
}

func auth(router *gin.RouterGroup, handler *authCtrl.Controller) {
	router.POST("/login", handler.Login)
	router.POST("/logout", handler.Logout)
	router.POST("/register", handler.Register)
	router.GET("/verify", handler.VerifyEmail)
	router.POST("/resend", handler.Resend)
}

func predict(router *gin.RouterGroup, handler *predictCtrl.Controller) {
	predict := router.Group("predict")

	predict.POST("/", handler.SubmitData)
}
