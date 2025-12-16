package handler

import (
	"database/sql"
	"lab1/internal/app/ds"
	"lab1/internal/app/role"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"bytes"
	"encoding/json"
)

type FormattedGrowthRequest struct {
	ID          uint    `json:"ID"`
	Status      string  `json:"Status"`
	DateCreate  string  `json:"DateCreate"`
	Creator     string  `json:"Creator"`
	DateUpdate  string  `json:"DateUpdate"`
	DateFinish  string  `json:"DateFinish"`
	Moderator   string  `json:"Moderator"`
	CurData     int     `json:"CurData"`
	StartPeriod string  `json:"StartPeriod"`
	EndPeriod   string  `json:"EndPeriod"`
	Result      float64 `json:"Result"`
}

// formatTime — вспомогательная функция
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02.01.06")
}

// formatNullTime — для sql.NullTime
func formatNullTime(nt interface{}) string {
	switch v := nt.(type) {
	case time.Time:
		if v.IsZero() {
			return ""
		}
		return v.Format("02.01.06")
	case *time.Time:
		if v == nil || v.IsZero() {
			return ""
		}
		return v.Format("02.01.06")
	default:
		return ""
	}
}

func FormatRequest(req ds.GrowthRequest) FormattedGrowthRequest {
	var moderatorLogin string
	if req.Moderator != nil {
		moderatorLogin = req.Moderator.Login
	}

	return FormattedGrowthRequest{
		ID:          req.ID,
		Status:      req.Status,
		DateCreate:  formatTime(req.DateCreate),
		Creator:     req.Creator.Login,
		DateUpdate:  formatTime(req.DateUpdate),
		DateFinish:  formatNullTime(req.DateFinish.Time),
		Moderator:   moderatorLogin,
		CurData:     req.CurData,
		StartPeriod: formatTime(req.StartPeriod),
		EndPeriod:   formatTime(req.EndPeriod),
		Result:      req.Result,
	}
}

// GetCartInfo godoc
// @Summary Получение информации о корзине пользователя
// @Description Возвращает данные о корзине в зависимости от роли (creator или moderator)
// @Tags growth_requests
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 403 "Forbidden"
// @Failure 500 {object} map[string]string
// @Router /api/growth-requests/cart [get]
func (h *Handler) GetCartInfo(ctx *gin.Context) {
	claims, err := h.ParseAndValidateJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"growth_request_id": 0,
			"service_count":     0,
		})
		return
	}

	userRole := claims.Role
	userID := claims.UserID

	switch role.Role(userRole) {
	case role.Moderator:
		carts, err := h.Repository.GetAllCartInfo()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"carts": carts,
		})
		return

	case role.Creator:
		cartID, count, err := h.Repository.GetCartInfo(userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"growth_request_id": cartID,
			"service_count":     count,
		})
		return

	default:
		// ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
}

// GetGrowthRequests godoc
// @Summary Получить список заявок на рост 2
// @Description Возвращает все заявки с фильтрацией по статусу и дате
// @Tags growth_requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Статус заявки"
// @Param start_date query string false "Дата начала периода (формат: 02.01.06)"
// @Param end_date query string false "Дата конца периода (формат: 02.01.06)"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 "Forbidden"
// @Failure 500 {object} map[string]string
// @Router /api/growth-requests [get]
func (h *Handler) GetGrowthRequests(ctx *gin.Context) {
	userRole := ctx.GetInt("role")
	userID := GetUserID(ctx)

	status := ctx.Query("status")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("02.01.06", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
			return
		}
	}
	if endDateStr != "" {
		endDate, err = time.Parse("02.01.06", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
			return
		}
	}

	requests, err := h.Repository.GetGrowthRequests(status, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, req := range requests {
		if dateCreate, ok := req["date_create"].(time.Time); ok {
			req["date_create"] = formatTime(dateCreate)
		}
		if dateFinish, ok := req["date_finish"].(time.Time); ok {
			req["date_finish"] = formatTime(dateFinish)
		}
		if dateUpdate, ok := req["date_update"].(time.Time); ok {
			req["date_update"] = formatTime(dateUpdate)
		}
	}

	switch role.Role(userRole) {
	case role.Moderator:
		ctx.JSON(http.StatusOK, requests)
		return
	case role.Creator:
		creatorRequests := make([]map[string]interface{}, 0)
		for _, req := range requests {
			if req["creator_id"] == userID {
				creatorRequests = append(creatorRequests, req)
			}
		}
		ctx.JSON(http.StatusOK, creatorRequests)
		return
	default:
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
}


// GetGrowthRequestByID godoc
// @Summary Получить заявку по ID
// @Description Возвращает заявку и связанные с ней факторы
// @Tags growth_requests
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/growth-requests/{id} [get]
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

	ctx.JSON(http.StatusOK, gin.H{"growth_request": FormatRequest(req), "factors": factors})
}

// UpdateGrowthRequest godoc
// @Summary Обновить заявку на рост
// @Description Обновляет данные заявки (текущие данные, период)
// @Tags growth_requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param input body object true "Поля для обновления"
// @Success 200 {object} FormattedGrowthRequest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/growth-requests/{id} [put]
func (h *Handler) UpdateGrowthRequest(ctx *gin.Context) {
	type updateGrowthRequestInput struct {
		CurData     int    `json:"cur_data"`
		StartPeriod string `json:"start_period"`
		EndPeriod   string `json:"end_period"`
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
	updated, err := h.Repository.UpdateGrowthRequest(uint(id), input.CurData, startPeriod, endPeriod, GetUserID(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, FormatRequest(*updated))
}

// FormGrowthRequest godoc
// @Summary Сформировать заявку на рост
// @Description Изменяет статус черновика на "сформирован" при наличии всех обязательных полей
// @Tags growth_requests
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/growth-requests/{id}/form [put]
func (h *Handler) FormGrowthRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid growth request id"})
		return
	}

	creatorID := GetUserID(ctx)
	gr, err := h.Repository.GetGrowthRequestByIDAndCreator(uint(id), creatorID)
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
	gr.CreatorID = creatorID
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

// CompleteOrRejectGrowthRequest godoc
// @Summary Завершить или отклонить заявку
// @Description Меняет статус заявки на завершен или отклонен
// @Tags growth_requests
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param action query string true "Действие (complete или reject)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/growth-requests/{id}/complete [put]
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

	moderatorID := GetUserID(ctx)

	switch action {
	case "complete":
		// var sum float64
		// for _, f := range factors {
		// 	sum += f.DataGrowthFactor.Coeff * f.FactorNum
		// }
		// duration := growthRequest.EndPeriod.Sub(growthRequest.StartPeriod).Hours() / 24
		// growthRequest.Result = float64(growthRequest.CurData) + sum*duration
		// growthRequest.Status = "завершен"
		requestBody := map[string]interface{}{
			"moderator_id": &moderatorID,
			"growth_request": growthRequest,
			"factors": factors,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot encode body"})
			return
		}
		url := "http://0.0.0.0:8000/"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "cannot reach remote service"})
			return
		}
		defer resp.Body.Close()
	case "reject":
		now := time.Now()
		growthRequest.Status = "отклонен"
		growthRequest.ModeratorID = &moderatorID
		growthRequest.DateFinish = sql.NullTime{Time: now, Valid: true}
		growthRequest.DateUpdate = now

		if err := h.Repository.SaveGrowthRequest(growthRequest); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}


// UpdateGrowthRequestResult godoc
// @Summary Добавить результат в заявку
// @Description Добавляет в заявку результат
// @Tags growth_requests
// @Produce json
// @Param id path int true "ID заявки"
// @Param input body object true "Поля для обновления"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/growth-requests/{id}/result [put]
func (h *Handler) UpdateGrowthRequestResult(ctx *gin.Context) {
	const defToken = "ABCDEF12"

	type updateGrowthRequestResultInput struct {
		Result     float64    `json:"result"`
		Token      string     `json:"token"`
		ModeratorID uint       `json:"moderator_id"`
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input updateGrowthRequestResultInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Token != defToken {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	growthRequest, _, err := h.Repository.GetGrowthRequestByIDWithFactors(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	now := time.Now()
	growthRequest.Status = "завершён"
	growthRequest.ModeratorID = &input.ModeratorID
	growthRequest.DateFinish = sql.NullTime{Time: now, Valid: true}
	growthRequest.DateUpdate = now
	growthRequest.Result = input.Result

	if err := h.Repository.SaveGrowthRequest(growthRequest); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}



// DeleteGrowthRequest godoc
// @Summary Удалить заявку на рост
// @Description Удаляет заявку пользователя
// @Tags growth_requests
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/growth-requests/{id} [delete]
func (h *Handler) DeleteGrowthRequest(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.Repository.DeleteGrowthRequest(uint(id), GetUserID(ctx))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "growth request deleted successfully"})
}
