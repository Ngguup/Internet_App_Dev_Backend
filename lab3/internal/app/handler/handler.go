package handler

import (
	"lab1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/dataGrowthHome", h.GetAllDataGrowthFactors_)
	router.GET("/dataGrowthHome/:id", h.GetDataGrowthFactorById_)
	router.GET("/growthRequest/:id", h.GrowthRequest_)
	// router.POST("/delete-dataGrowthFactor", h.DeleteDataGrowthFactor_)
	router.POST("/add-dataGrowthFactor", h.AddDataGrowthFactor_)
	router.POST("/delete-growthRequest", h.DeleteDataGrowthFactor_)

	api := router.Group("/api")
	{
		api.GET("/data-growth-factors", h.GetAllDataGrowthFactors)
		api.GET("/data-growth-factors/:id", h.GetDataGrowthFactorByID)
		api.POST("/data-growth-factors", h.CreateDataGrowthFactor)
		api.PUT("/data-growth-factors/:id", h.UpdateDataGrowthFactor)
		api.DELETE("/data-growth-factors/:id", h.DeleteDataGrowthFactor)
		api.POST("/add-data-growth-factor-to-draft/:id", h.AddDataGrowthFactorToDraft)
		api.POST("/data-growth-factors/:id/image", h.UploadDataGrowthFactorImage)

		api.GET("/growth-requests-cart", h.GetCartInfo)
		api.GET("/growth-requests", h.GetGrowthRequests)
		api.GET("/growth-requests/:id", h.GetGrowthRequestByID)
		api.PUT("/growth-requests/:id", h.UpdateGrowthRequest)
		api.PUT("/growth-requests/:id/form", h.FormGrowthRequest)
		api.PUT("/growth-requests/:id/complete", h.CompleteOrRejectGrowthRequest)
		api.DELETE("/growth-requests/:id", h.DeleteGrowthRequest)

		api.DELETE("/growth_request_data_growth_factors/:id", h.DeleteDataGrowthFactorFromDraft)
		api.PUT("/growth_request_data_growth_factors/:id", h.UpdateFactorNum)

		api.POST("/users/register", h.RegisterUser)
		api.GET("/users/me", h.GetCurrentUser)
		api.PUT("/users/me", h.UpdateUser)
	}
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
