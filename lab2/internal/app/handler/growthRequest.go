package handler
import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)


func (h *Handler) DeleteDataGrowthFactor(ctx *gin.Context) {
	err := h.Repository.DeleteGrowthRequest()
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return
	}

	ctx.Redirect(http.StatusFound, "/dataGrowthHome")
}