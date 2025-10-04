package handler

import (
	"lab1/internal/app/repository"
	"net/http"
	"strconv"

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

func (h *Handler) GrowthRequest(ctx *gin.Context) {
	var err error

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	growthRequest, factorNums, err := h.Repository.GetGrowthRequestByID(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "growthRequest.html", gin.H{
		"growthRequest":     growthRequest,
		"dataGrowthFactors": growthRequest.Components,
		"factorNums":        factorNums,
		"startPeriod":       growthRequest.StartPeriod,
		"endPeriod":         growthRequest.EndPeriod,
	})
}

func (h *Handler) GetDataGrowthFactors(ctx *gin.Context) {
	var dataGrowthFactors []repository.DataGrowthFactor
	var err error
	var growthRequest repository.GrowthRequest
	growthRequestID :=2
	growthRequestID -= 1

	growthRequest, _, err = h.Repository.GetGrowthRequestByID(growthRequestID)
	if err != nil {
		logrus.Error(err)
	}

	growthFactorSearchQuery := ctx.Query("query")
	if growthFactorSearchQuery == "" {
		dataGrowthFactors, err = h.Repository.GetDataGrowthFactors()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		dataGrowthFactors, err = h.Repository.GetDataGrowthFactorsByTitle(growthFactorSearchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"numOfDataGrowthFactors": len(growthRequest.Components),
		"dataGrowthFactors":      dataGrowthFactors,
		"query":                  growthFactorSearchQuery,
		"growthRequestID": growthRequestID,
	})
}

func (h *Handler) GetDataGrowthFactor(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	dataGrowthFactor, err := h.Repository.GetDataGrowthFactor(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "dataGrowthFactor.html", gin.H{
		"dataGrowthFactor": dataGrowthFactor,
	})
}
