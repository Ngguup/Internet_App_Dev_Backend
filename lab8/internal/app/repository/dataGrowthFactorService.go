package repository

import (
	"lab1/internal/app/ds"
	"fmt"
	"errors"
	"gorm.io/gorm"
	"time"

	"context"
	"strings"
	"lab1/internal/app/dsn"
	"github.com/minio/minio-go/v7"
)

func (r *Repository) GetAllDataGrowthFactors(title string, minCoeff, maxCoeff *float64) ([]ds.DataGrowthFactor, error) {
	var factors []ds.DataGrowthFactor

	query := r.db.Where("is_delete = false")

	if title != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+title+"%")
	}

	if minCoeff != nil {
		query = query.Where("coeff >= ?", *minCoeff)
	}
	if maxCoeff != nil {
		query = query.Where("coeff <= ?", *maxCoeff)
	}

	err := query.Find(&factors).Error
	return factors, err
}


func (r *Repository) GetDataGrowthFactorByID(id uint) (*ds.DataGrowthFactor, error) {
	var factor ds.DataGrowthFactor
	err := r.db.Where("id = ? AND is_delete = false", id).First(&factor).Error
	if err != nil {
		return nil, err
	}
	return &factor, nil
}

func (r *Repository) CreateDataGrowthFactor(factor *ds.DataGrowthFactor) error {
	return r.db.Create(factor).Error
}

func (r *Repository) UpdateDataGrowthFactor(id uint, input *ds.DataGrowthFactor) error {
	return r.db.Model(&ds.DataGrowthFactor{}).
		Where("id = ? AND is_delete = false", id).
		Updates(map[string]interface{}{
			"title":       input.Title,
			"coeff":       input.Coeff,
			"description": input.Description,
			"attribute":   input.Attribute,
		}).Error
}

func (r *Repository) DeleteDataGrowthFactor(id uint) error {
	var factor ds.DataGrowthFactor

	if err := r.db.Unscoped().Where("id = ?", id).First(&factor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	if factor.IsDelete {
		return fmt.Errorf("data_growth_factor already deleted")
	}

	if factor.Image != "" {
		client, bucketName, err := dsn.GetMinioClient()
		if err != nil {
			return fmt.Errorf("failed to connect to MinIO: %w", err)
		}

		imageURL := factor.Image
		parts := strings.Split(imageURL, "/")
		objectName := parts[len(parts)-1] 

		err = client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete image from MinIO: %w", err)
		}
	}

	if err := r.db.Model(&factor).Update("is_delete", true).Error; err != nil {
		return fmt.Errorf("failed to mark deleted: %w", err)
	}

	return nil
}


func (r *Repository) AddDataGrowthFactorToDraft(dataGrowthFactorID uint, creatorID uint) error {
    var draft ds.GrowthRequest

    err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&draft).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        draft = ds.GrowthRequest{
            CreatorID:  creatorID,
            Status:     "черновик",
            DateCreate: time.Now(),
        }
        if err := r.db.Create(&draft).Error; err != nil {
            return err
        }
    } else if err != nil {
        return err
    }

    grdf := ds.GrowthRequestDataGrowthFactor{
        GrowthRequestID:    draft.ID,
        DataGrowthFactorID: dataGrowthFactorID,
        FactorNum:          1.0, 
    }

    if err := r.db.Create(&grdf).Error; err != nil {
        return err
    }

    return nil
}

func (r *Repository) UpdateDataGrowthFactorImage(id uint, imageURL string) error {
    return r.db.Model(&ds.DataGrowthFactor{}).
        Where("id = ?", id).
        Update("image", imageURL).Error
}


