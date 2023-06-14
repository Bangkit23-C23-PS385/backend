package v1

import (
	authCtrl "backend/src/controller/v1/auth"
	predictCtrl "backend/src/controller/v1/predict"
	profileCtrl "backend/src/controller/v1/profile"
	"backend/src/database"
	authRepo "backend/src/repository/v1/auth"
	predictRepo "backend/src/repository/v1/predict"
	profileRepo "backend/src/repository/v1/profile"
	verifRepo "backend/src/repository/v1/verification"
	authSvc "backend/src/service/v1/auth"
	predictSvc "backend/src/service/v1/predict"
	profileSvc "backend/src/service/v1/profile"
	verifSvc "backend/src/service/v1/verification"

	"github.com/gin-gonic/gin"
)

func RouteLoader(router *gin.Engine, db database.DB) {
	predictRepository := predictRepo.NewRepository(db)
	verifRepository := verifRepo.NewRepository(db)
	authRepository := authRepo.NewRepository(db)
	profileRepository := profileRepo.NewRepository(db)

	predictService := predictSvc.NewService(predictRepository)
	verifService := verifSvc.NewService(verifRepository)
	authService := authSvc.NewService(verifService, authRepository)
	profileService := profileSvc.NewService(profileRepository)

	predictController := predictCtrl.NewController(predictService)
	authController := authCtrl.NewController(authService)
	profileController := profileCtrl.NewController(profileService)

	v1 := router.Group("v1")

	auth(v1, authController)
	predict(v1, predictController)
	profile(v1, profileController)
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

	symptoms := router.Group("symptoms")

	symptoms.GET("", handler.GetSymptoms)
}

func profile(router *gin.RouterGroup, handler *profileCtrl.Controller) {
	router.GET("/profile", handler.GetProfile)
	router.POST("/profile", handler.CreateProfile)
	router.PUT("/profile", handler.UpdateProfile)
	router.DELETE("/profile", handler.DeleteProfile)
}
