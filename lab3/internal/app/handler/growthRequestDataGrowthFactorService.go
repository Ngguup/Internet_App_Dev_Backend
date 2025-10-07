package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteDataGrowthFactorFromDraft(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid data_growth_factor_id"})
		return
	}


	err = h.Repository.DeleteFromDraft(uint(id), CreatorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "DataGrowthFactor removed from draft successfully"})
}

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

	creatorID := 1 

	if err := h.Repository.UpdateFactorNum(uint(id), body.FactorNum, uint(creatorID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "factor_num updated successfully"})
}

