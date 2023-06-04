package v1

import (
	authCtrl "backend/src/controller/v1/auth"
	"backend/src/database"
	authRepo "backend/src/repository/v1/auth"
	verifRepo "backend/src/repository/v1/verification"
	authSvc "backend/src/service/v1/auth"
	verifSvc "backend/src/service/v1/verification"

	"github.com/gin-gonic/gin"
)

func RouteLoader(router *gin.Engine, db database.DB) {
	verifRepository := verifRepo.NewRepository(db)
	authRepository := authRepo.NewRepository(db)

	verifService := verifSvc.NewService(verifRepository)
	authService := authSvc.NewService(verifService, authRepository)

	authController := authCtrl.NewController(authService)

	v1 := router.Group("v1")

	auth(v1, authController)
}

func auth(router *gin.RouterGroup, handler *authCtrl.Controller) {
	router.POST("/login", handler.Login)
	router.POST("/logout", handler.Logout)
	router.POST("/register", handler.Register)
	router.GET("/verify", handler.VerifyEmail)
	router.POST("/resend", handler.Resend)
}
