package handler

import (
	"context"
	"errors"
	"fmt"
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// GetAllDataGrowthFactors godoc
// @Summary      Получить список услуг
// @Description  Возвращает все услуги (факторы роста) с возможностью фильтрации по названию
// @Tags         data_growth_factors
// @Produce      json
// @Param        title  query  string  false  "Фильтр по названию2"
// @Param        min_coeff  query  string  false "Минимальный коэффициент" 
// @Param        max_coeff  query  string  false "Максимальный коэффициент"
// @Success      200  {array}  ds.DataGrowthFactor
// @Router       /api/data-growth-factors [get]
func (h *Handler) GetAllDataGrowthFactors(ctx *gin.Context) {
	title := ctx.Query("title")

	// Читаем диапазон коэффициента
	minCoeffStr := ctx.Query("min_coeff")
	maxCoeffStr := ctx.Query("max_coeff")

	var minCoeff, maxCoeff *float64

	if minCoeffStr != "" {
		val, err := strconv.ParseFloat(minCoeffStr, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_coeff"})
			return
		}
		minCoeff = &val
	}

	if maxCoeffStr != "" {
		val, err := strconv.ParseFloat(maxCoeffStr, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_coeff"})
			return
		}
		maxCoeff = &val
	}

	factors, err := h.Repository.GetAllDataGrowthFactors(title, minCoeff, maxCoeff)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, factors)
}


// GetDataGrowthFactorByID godoc
// @Summary      Получить услугу по ID
// @Description  Возвращает услугу (фактор роста) по её идентификатору
// @Tags         data_growth_factors
// @Produce      json
// @Param        id   path  int  true  "ID услуги"
// @Success      200  {object}  ds.DataGrowthFactor
// @Failure      404  {object}  map[string]string  "Услуга не найдена"
// @Router       /api/data-growth-factors/{id} [get]
func (h *Handler) GetDataGrowthFactorByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	factor, err := h.Repository.GetDataGrowthFactorByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if factor == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Услуга не найдена"})
		return
	}

	ctx.JSON(http.StatusOK, factor)
}

// CreateDataGrowthFactor godoc
// @Summary      Добавить услугу
// @Description  Создаёт новую услугу (фактор роста)
// @Tags         data_growth_factors
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        data  body  ds.DataGrowthFactor  true  "Данные новой услуги"
// @Success      201  {object}  ds.DataGrowthFactor
// @Router       /api/data-growth-factors [post]
func (h *Handler) CreateDataGrowthFactor(ctx *gin.Context) {
	var input ds.DataGrowthFactor
	if err := ctx.BindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	newFactor := ds.DataGrowthFactor{
		Title:       input.Title,
		Description: input.Description,
		Attribute:   input.Attribute,
		Coeff:       input.Coeff,
		Image:       "",
		IsDelete:    false,
	}

	err := h.Repository.CreateDataGrowthFactor(&newFactor)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, newFactor)
}

// UpdateDataGrowthFactor godoc
// @Summary      Обновить услугу
// @Description  Обновляет информацию об услуге (факторе роста)
// @Tags         data_growth_factors
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        id    path  int  true  "ID услуги"
// @Param        data  body  ds.DataGrowthFactor  true  "Обновлённые данные"
// @Success      200   {object}  map[string]interface{}
// @Router       /api/data-growth-factors/{id} [put]
func (h *Handler) UpdateDataGrowthFactor(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var input ds.DataGrowthFactor
	if err := ctx.BindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateDataGrowthFactor(uint(id), &input)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "updated", "time": time.Now()})
}

// DeleteDataGrowthFactor godoc
// @Summary      Удалить услугу
// @Description  Удаляет услугу (меняет флаг удаления или полностью удаляет запись)
// @Tags         data_growth_factors
// @Security BearerAuth
// @Param        id   path  int  true  "ID услуги"
// @Success      200  {object}  map[string]string  "DataGrowthFactor deleted successfully"
// @Failure      404  {object}  map[string]string  "data_growth_factor not found"
// @Router       /api/data-growth-factors/{id} [delete]
func (h *Handler) DeleteDataGrowthFactor(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.Repository.DeleteDataGrowthFactor(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "data_growth_factor not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "DataGrowthFactor deleted successfully"})
}

// AddDataGrowthFactorToDraft godoc
// @Summary      Добавить услугу в заявку-черновик
// @Description  Добавляет выбранную услугу в текущую заявку пользователя со статусом "черновик"
// @Tags         data_growth_factors
// @Security BearerAuth
// @Param        id   path  int  true  "ID услуги"
// @Success      200  {object}  map[string]string  "DataGrowthFactor added to draft successfully"
// @Router       /api/data-growth-factors/{id}/add [post]
func (h *Handler) AddDataGrowthFactorToDraft(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.Repository.AddDataGrowthFactorToDraft(uint(id), GetUserID(ctx)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "DataGrowthFactor added to draft successfully"})
}

// UploadDataGrowthFactorImage godoc
// @Summary      Загрузить изображение услуги
// @Description  Загружает изображение в MinIO и сохраняет URL в базе данных
// @Tags         data_growth_factors
// @Accept       multipart/form-data
// @Produce      json
// @Security BearerAuth
// @Param        id     path      int     true   "ID услуги"
// @Param        image  formData  file    true   "Изображение"
// @Success      200    {object}  map[string]string  "image_url"
// @Router       /api/data-growth-factors/{id}/image [post]
func (h *Handler) UploadDataGrowthFactorImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file not provided"})
		return
	}
	defer file.Close()

	client, bucketName, err := dsn.GetMinioClient()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot connect to minio"})
		return
	}

	imageName := header.Filename

	_, err = client.PutObject(
		context.Background(),
		bucketName,
		imageName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot upload image"})
		return
	}

	imageURL := fmt.Sprintf("http://%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), bucketName, imageName)
	h.Repository.UpdateDataGrowthFactorImage(uint(id), imageURL)

	ctx.JSON(http.StatusOK, gin.H{"image_url": imageURL})
}
