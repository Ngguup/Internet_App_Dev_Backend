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

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/dataGrowthHome", h.GetAllDataGrowthFactors)
	router.GET("/dataGrowthHome/:id", h.GetDataGrowthFactorById)
	router.GET("/growthRequest/:id", h.GrowthRequest)
	// router.POST("/delete-dataGrowthFactor", h.DeleteDataGrowthFactor)
	router.POST("/add-dataGrowthFactor", h.AddDataGrowthFactor)
	router.POST("/delete-growthRequest", h.DeleteDataGrowthFactor)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
