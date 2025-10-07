package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/sirupsen/logrus"
	"database/sql"
)

const CreatorID uint = 1 
const ModeratorID uint = 2

func (h *Handler) GetCartInfo(ctx *gin.Context) {
	cartID, count, err := h.Repository.GetCartInfo(CreatorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"growth_request_id": cartID,
		"service_count":     count,
	})
}

func (h *Handler) GetGrowthRequests(ctx *gin.Context) {
	type filterInput struct {
		Status    string `json:"status"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	var input filterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var startDate, endDate time.Time
	var err error

	if input.StartDate != "" {
		startDate, err = time.Parse("02.01.06", input.StartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
			return
		}
	}
	if input.EndDate != "" {
		endDate, err = time.Parse("02.01.06", input.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
			return
		}
	}

	requests, err := h.Repository.GetGrowthRequests(input.Status, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, requests)
}


func (h *Handler) GetGrowthRequestByID(ctx *gin.Context) {
	id := ctx.Param("id")

	req, factors, err := h.Repository.GetGrowthRequestByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if req.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "growth request not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"growth_request": req,
		"factors":        factors,
	})
}

func (h *Handler) UpdateGrowthRequest(ctx *gin.Context) {
	type updateGrowthRequestInput struct {
		CurData     int       `json:"cur_data"`
		StartPeriod string    `json:"start_period"`
		EndPeriod   string    `json:"end_period"`
	}
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input updateGrowthRequestInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	startPeriod, err := time.Parse("02.01.06", input.StartPeriod)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_period"})
		return
	}
	endPeriod, err := time.Parse("02.01.06", input.EndPeriod)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_period"})
		return
	}
	updated, err := h.Repository.UpdateGrowthRequest(uint(id), input.CurData, startPeriod, endPeriod)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

func (h *Handler) FormGrowthRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid growth request id"})
		return
	}

	gr, err := h.Repository.GetGrowthRequestByIDAndCreator(uint(id), CreatorID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "growth request not found"})
		return
	}

	if gr.Status != "черновик" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only draft requests can be formed"})
		return
	}

	if gr.CurData == 0 || gr.StartPeriod.IsZero() || gr.EndPeriod.IsZero() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "mandatory fields are missing"})
		return
	}

	gr.Status = "сформирован"
	gr.DateCreate = time.Now()
	gr.CreatorID = CreatorID
	gr.DateUpdate = time.Now()

	if err := h.Repository.UpdateGrowthRequestForm(gr); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot form growth request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           gr.ID,
		"status":       gr.Status,
		"cur_data":     gr.CurData,
		"start_period": gr.StartPeriod.Format("02.01.06"),
		"end_period":   gr.EndPeriod.Format("02.01.06"),
		"date_create":  gr.DateCreate.Format("02.01.06"),
		"creator_id":   gr.CreatorID,
	})
}

func (h *Handler) CompleteOrRejectGrowthRequest(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    action := ctx.Query("action") 
    if action != "complete" && action != "reject" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
        return
    }

    growthRequest, factors, err := h.Repository.GetGrowthRequestByIDWithFactors(uint(id))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if growthRequest.Status != "сформирован" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "only formed requests can be completed or rejected"})
        return
    }

    now := time.Now()

    switch action {
		case "complete":
			var sum float64
			for _, f := range factors {
				sum += f.DataGrowthFactor.Coeff * f.FactorNum
			}
			duration := growthRequest.EndPeriod.Sub(growthRequest.StartPeriod).Hours() / 24
			growthRequest.Result = float64(growthRequest.CurData) + sum*duration
			growthRequest.Status = "завершен"
		case "reject":
			growthRequest.Status = "отклонен"
		}

    growthRequest.ModeratorID = ModeratorID
    growthRequest.DateFinish = sql.NullTime{Time: now, Valid: true}
	growthRequest.DateUpdate = now

    if err := h.Repository.SaveGrowthRequest(growthRequest); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "success", "result": growthRequest.Result})
}

func (h *Handler) DeleteGrowthRequest(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.Repository.DeleteGrowthRequest(uint(id), uint(CreatorID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "growth request deleted successfully"})
}