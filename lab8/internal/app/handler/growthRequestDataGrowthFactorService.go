package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteDataGrowthFactorFromDraft godoc
// @Summary Удалить фактор роста из черновика
// @Description Удаляет DataGrowthFactor из черновика текущего пользователя
// @Tags growth_request_data_growth_factors
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID фактора роста"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/growth-request-data-growth-factors/{id} [delete]
func (h *Handler) DeleteDataGrowthFactorFromDraft(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid data_growth_factor_id"})
		return
	}

	err = h.Repository.DeleteFromDraft(uint(id), GetUserID(ctx))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "DataGrowthFactor removed from draft successfully"})
}

// UpdateFactorNum godoc
// @Summary Обновить значение factor_num
// @Description Обновляет числовое значение factor_num для DataGrowthFactor в корзине текущего пользователя
// @Tags growth_request_data_growth_factors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID фактора роста"
// @Param body body object true "Новые данные"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/growth-request-data-growth-factors/{id} [put]
func (h *Handler) UpdateFactorNum(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid data_growth_factor_id"})
		return
	}

	var body struct {
		FactorNum float64 `json:"factor_num"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := h.Repository.UpdateFactorNum(uint(id), body.FactorNum, GetUserID(ctx)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "factor_num updated successfully"})
}
