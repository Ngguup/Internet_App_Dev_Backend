package handler

import (
	"lab1/internal/app/ds"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUser(ctx *gin.Context) {
	var input struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user := ds.Users{
		Login:    input.Login,
		Password: input.Password,
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *Handler) GetCurrentUser(ctx *gin.Context) {
	const creatorID = 1

	user, err := h.Repository.GetUserByID(uint(creatorID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"login":        user.Login,
		"is_moderator": user.IsModerator,
	})
}

func (h *Handler) UpdateUser(ctx *gin.Context) {
	const creatorID = 1

	var input struct {
		Login       string `json:"login"`
		Password    string `json:"password"`
		IsModerator *bool  `json:"is_moderator"` 
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.Repository.UpdateUser(uint(creatorID), input.Login, input.Password, input.IsModerator); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
