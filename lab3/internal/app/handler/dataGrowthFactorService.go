package handler

import (
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn"
	"net/http"
	"strconv"
	"time"
	"errors"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"context"
    "fmt"
    "os"
    "github.com/minio/minio-go/v7"
)

func (h *Handler) GetAllDataGrowthFactors(ctx *gin.Context) {
	title := ctx.Query("title")

	factors, err := h.Repository.GetAllDataGrowthFactors(title)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, factors)
}

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

func (h *Handler) AddDataGrowthFactorToDraft(ctx *gin.Context) {
    idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

    if err := h.Repository.AddDataGrowthFactorToDraft(uint(id)); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "DataGrowthFactor added to draft successfully"})
}

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



