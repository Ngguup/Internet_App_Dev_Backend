package handler

import (
	"net/http"
	"strconv"
	"strings"

	"lab1/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GrowthRequest_(ctx *gin.Context) {
	var growthRequest ds.GrowthRequest
	var dataGrowthFactors []ds.DataGrowthFactor
	var factorNums []float64
	var err error

	var startPeriod string
	var endPeriod string

	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		// ctx.String(http.StatusNotFound, "404 page not found")
		ctx.Redirect(http.StatusFound, "/dataGrowthHome")
		return
	}

	growthRequest, dataGrowthFactors, factorNums, err = h.Repository.GetGrowthRequestByID_(id)
	if err != nil {
		logrus.Error(err)
		// ctx.String(http.StatusNotFound, "404 page not found")
		ctx.Redirect(http.StatusFound, "/dataGrowthHome")
		return
	}

	//Расчёт роста данных:
	if growthRequest.Status == "завершён" {
		days := growthRequest.EndPeriod.Sub(growthRequest.StartPeriod).Hours() / 24
		if days < 0 {
			days = 0
		}

		var dataForecast float64
		for i := range factorNums {
			dataForecast += factorNums[i] * dataGrowthFactors[i].Coeff * days
		}
	}
	//

	startPeriod = growthRequest.StartPeriod.Format("02.01.06")
	endPeriod = growthRequest.EndPeriod.Format("02.01.06")

	ctx.HTML(http.StatusOK, "growthRequest.page.tmpl", gin.H{
		"growthRequest":     growthRequest,
		"dataGrowthFactors": dataGrowthFactors,
		"factorNums":        factorNums,
		"startPeriod":       startPeriod,
		"endPeriod":         endPeriod,
	})
}

func (h *Handler) GetAllDataGrowthFactors_(ctx *gin.Context) {
	var dataGrowthFactors []ds.DataGrowthFactor
	var err error

	query := ctx.Query("query")
	if query == "" {
		dataGrowthFactors, err = h.Repository.GetAllDataGrowthFactors_()
	} else {
		dataGrowthFactors, err = h.Repository.SearchDataGrowthFactorsByName_(query)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	growthRequestID, numOfDataGrowthFactors := h.Repository.GetGrowthRequestData_()

	ctx.HTML(http.StatusOK, "index.page.tmpl", gin.H{
		"dataGrowthFactors":      dataGrowthFactors,
		"numOfDataGrowthFactors": numOfDataGrowthFactors,
		"growthRequestID":        growthRequestID,
		"query":                  query,
	})
}

func (h *Handler) GetDataGrowthFactorById_(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	dataGrowthFactor, err := h.Repository.GetDataGrowthFactorById_(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "dataGrowthFactor.page.tmpl", dataGrowthFactor)
}

func (h *Handler) AddDataGrowthFactor_(ctx *gin.Context) {
	strId := ctx.PostForm("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	err = h.Repository.AddDataGrowthFactor_(uint(id))
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return
	}

	ctx.Redirect(http.StatusFound, "/dataGrowthHome")
}

// func (h *Handler) DeleteDataGrowthFactor_(ctx *gin.Context) {
// 	// считываем значение из формы, которую мы добавим в наш шаблон
// 	strId := ctx.PostForm("data_growth_factor_id")
// 	id, err := strconv.Atoi(strId)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 	}
// 	// Вызов функции добавления чата в заявку
// 	err = h.Repository.DeleteDataGrowthFactor_(uint(id))
// 	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
// 		return
// 	}

// 	// после вызова сразу произойдет обновление страницы
// 	ctx.Redirect(http.StatusFound, "/dataGrowthHome")
// }
