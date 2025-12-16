package handler

import (
	"lab1/internal/app/config"
	"lab1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"lab1/internal/app/redis"
	"lab1/internal/app/role"

	_ "lab1/docs"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"     
	ginSwagger "github.com/swaggo/gin-swagger" 
)

type Handler struct {
	Repository *repository.Repository
	Config     *config.Config
	Redis      *redis.Client
}

func NewHandler(rep *repository.Repository, conf *config.Config, redis *redis.Client) *Handler {
	return &Handler{
		Repository: rep,
		Config:     conf,
		Redis:      redis,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/dataGrowthHome", h.GetAllDataGrowthFactors_)
	router.GET("/dataGrowthHome/:id", h.GetDataGrowthFactorById_)
	router.GET("/growthRequest/:id", h.GrowthRequest_)
	// router.POST("/delete-dataGrowthFactor", h.DeleteDataGrowthFactor_)
	router.POST("/add-dataGrowthFactor", h.AddDataGrowthFactor_)
	router.POST("/delete-growthRequest", h.DeleteDataGrowthFactor_)

	api := router.Group("/api")
	{
		api.GET("/growth-requests/cart", h.GetCartInfo)
		api.GET("/data-growth-factors", h.GetAllDataGrowthFactors)
		api.GET("/data-growth-factors/:id", h.GetDataGrowthFactorByID)
		// api.POST("/users/register", h.RegisterUser)
		api.POST("/users/register", h.Register)
		api.POST("/users/login", h.Login)

		api.PUT("growth-requests/:id/result", h.UpdateGrowthRequestResult)
	}

	apiCreatorModerator := router.Group("/api")
	apiCreatorModerator.Use(h.WithAuthCheck(role.Creator, role.Moderator))
	{
		apiCreatorModerator.POST("/data-growth-factors/:id/add", h.AddDataGrowthFactorToDraft)

		apiCreatorModerator.GET("/growth-requests", h.GetGrowthRequests)
		apiCreatorModerator.GET("/growth-requests/:id", h.GetGrowthRequestByID)

		apiCreatorModerator.PUT("/growth-requests/:id", h.UpdateGrowthRequest)
		apiCreatorModerator.PUT("/growth-requests/:id/form", h.FormGrowthRequest)
		apiCreatorModerator.DELETE("/growth-requests/:id", h.DeleteGrowthRequest)

		apiCreatorModerator.DELETE("/growth-request-data-growth-factors/:id", h.DeleteDataGrowthFactorFromDraft)
		apiCreatorModerator.PUT("/growth-request-data-growth-factors/:id", h.UpdateFactorNum)

		apiCreatorModerator.GET("/users/me", h.GetCurrentUser)
		apiCreatorModerator.PUT("/users/me", h.UpdateUser)
		apiCreatorModerator.POST("/users/logout", h.Logout)
	}

	apiModerator := router.Group("/api")
	apiModerator.Use(h.WithAuthCheck(role.Moderator))
	{
		apiModerator.POST("/data-growth-factors", h.CreateDataGrowthFactor)
		apiModerator.PUT("/data-growth-factors/:id", h.UpdateDataGrowthFactor)
		apiModerator.DELETE("/data-growth-factors/:id", h.DeleteDataGrowthFactor)
		apiModerator.POST("/data-growth-factors/:id/image", h.UploadDataGrowthFactorImage)

		apiModerator.PUT("/growth-requests/:id/complete", h.CompleteOrRejectGrowthRequest)
	}

	// router.Use(h.WithAuthCheck(role.Manager, role.Admin)).GET("/ping", h.Ping)
	// router.GET("/ping", h.WithAuthCheck(role.Manager, role.Admin), h.Ping)
	// authorized := router.Group("/")
	// authorized.Use(h.WithAuthCheck(role.Manager, role.Admin, role.Buyer))
	// {
	// 	authorized.GET("/ping", h.Ping)
	// }
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
