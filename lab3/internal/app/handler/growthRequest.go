package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteDataGrowthFactor_(ctx *gin.Context) {
	// Вызов функции добавления чата в заявку
	err := h.Repository.DeleteGrowthRequest_()
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return
	}

	// после вызова сразу произойдет обновление страницы
	ctx.Redirect(http.StatusFound, "/dataGrowthHome")
}
