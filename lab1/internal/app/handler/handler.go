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

func (h *Handler) Basket(ctx *gin.Context) {
	var basket []repository.Basket
	var err error
	var curBasket repository.Basket

	basket, err = h.Repository.GetBasket()
	if err != nil {
		logrus.Error(err)
	}

	for _, t := range basket { if t.Status { curBasket = t; break } }

	ctx.HTML(http.StatusOK, "basket.html", gin.H{
		"basket": curBasket,
	})
}

func (h *Handler) GetOrders(ctx *gin.Context) {
	var orders []repository.Order
	var err error
	var basket []repository.Basket
	var numOfOrders int

	basket, err = h.Repository.GetBasket()
	if err != nil {
		logrus.Error(err)
	}

	for _, t := range basket { if t.Status { numOfOrders = len(t.Components); break } }

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		orders, err = h.Repository.GetOrders()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		orders, err = h.Repository.GetOrdersByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"numOfOrders": numOfOrders,
		"orders":      orders,
		"query":       searchQuery,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr) 
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order": order,
	})
}
