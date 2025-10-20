package repository

import (
	"errors"
	"lab1/internal/app/ds"
	"time"

	"gorm.io/gorm"
	"fmt"

	"github.com/gin-gonic/gin"
)



func (r *Repository) GetCartInfo(creatorID uint) (uint, int64, error) {
	var req ds.GrowthRequest
	if err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").
		First(&req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	var count int64
	if err := r.db.Model(&ds.GrowthRequestDataGrowthFactor{}).
		Where("growth_request_id = ?", req.ID).
		Count(&count).Error; err != nil {
		return req.ID, 0, err
	}

	return req.ID, count, nil
}

func (r *Repository) GetAllCartInfo() ([]gin.H, error) {
	var requests []ds.GrowthRequest
	if err := r.db.Where("status = ?", "черновик").Find(&requests).Error; err != nil {
		return nil, err
	}

	carts := make([]gin.H, 0, len(requests))
	for _, req := range requests {
		var count int64
		if err := r.db.Model(&ds.GrowthRequestDataGrowthFactor{}).
			Where("growth_request_id = ?", req.ID).
			Count(&count).Error; err != nil {
			return nil, err
		}

		carts = append(carts, gin.H{
			"growth_request_id": req.ID,
			"service_count":     count,
			"creator_id":        req.CreatorID, 
		})
	}

	return carts, nil
}




func (r *Repository) GetGrowthRequests(status string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	query := r.db.Model(&ds.GrowthRequest{}).
		Select(`growth_requests.id,
                growth_requests.status,
                growth_requests.date_create,
                growth_requests.date_finish,
                growth_requests.creator_id,
                growth_requests.moderator_id`).
		Where("growth_requests.status != ?", "удалён")

	if status != "" {
		query = query.Where("growth_requests.status = ?", status)
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("growth_requests.date_create BETWEEN ? AND ?", startDate, endDate)
	}

	var result []map[string]interface{}
	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}


func (r *Repository) GetGrowthRequestByID(id string) (ds.GrowthRequest, []map[string]interface{}, error) {
	var req ds.GrowthRequest
	if err := r.db.Preload("Creator").Preload("Moderator").First(&req, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.GrowthRequest{}, nil, nil
		}
		return ds.GrowthRequest{}, nil, err
	}

	rows, err := r.db.
		Table("data_growth_factors AS dgf").
		Select(`
			dgf.id, 
			dgf.title, 
			dgf.image, 
			dgf.coeff, 
			dgf.description, 
			dgf.is_delete,
			grdf.factor_num
		`).
		Joins("JOIN growth_request_data_growth_factors grdf ON grdf.data_growth_factor_id = dgf.id").
		Where("grdf.growth_request_id = ?", req.ID).
		Where("dgf.is_delete = false").
		Rows()
	if err != nil {
		return req, nil, err
	}
	defer rows.Close()

	var factors []map[string]interface{}
	for rows.Next() {
		var (
			id          uint
			title       string
			image       string
			coeff       float64
			description string
			isDelete    bool
			factorNum   float64
		)
		if err := rows.Scan(&id, &title, &image, &coeff, &description, &isDelete, &factorNum); err != nil {
			return req, nil, err
		}

		factors = append(factors, map[string]interface{}{
			"ID":          id,
			"Title":       title,
			"Image":       image,
			"Coeff":       coeff,
			"Description": description,
			"IsDelete":    isDelete,
			"FactorNum":   factorNum,
		})
	}

	return req, factors, nil
}


func (r *Repository) UpdateGrowthRequest(id uint, curData int, startPeriod, endPeriod time.Time, creatorID uint) (*ds.GrowthRequest, error) {
	var gr ds.GrowthRequest

	if err := r.db.Where("id = ? AND status = ? AND creator_id = ?", id, "черновик", creatorID).
		First(&gr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("growth request not found or not editable")
		}
		return nil, err
	}

	gr.CurData = curData
	gr.StartPeriod = startPeriod
	gr.EndPeriod = endPeriod
	gr.DateUpdate = time.Now()


	if err := r.db.Save(&gr).Error; err != nil {
		return nil, err
	}

	return &gr, nil
}

func (r *Repository) GetGrowthRequestByIDAndCreator(id uint, creatorID uint) (*ds.GrowthRequest, error) {
	var gr ds.GrowthRequest
	if err := r.db.Where("id = ? AND creator_id = ?", id, creatorID).First(&gr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &gr, nil
}

func (r *Repository) UpdateGrowthRequestForm(gr *ds.GrowthRequest) error {
	return r.db.Save(gr).Error
}

func (r *Repository) GetGrowthRequestByIDWithFactors(id uint) (*ds.GrowthRequest, []ds.GrowthRequestDataGrowthFactor, error) {
    var gr ds.GrowthRequest
    if err := r.db.First(&gr, id).Error; err != nil {
        return nil, nil, err
    }

    var factors []ds.GrowthRequestDataGrowthFactor
    if err := r.db.Preload("DataGrowthFactor").Where("growth_request_id = ?", id).Find(&factors).Error; err != nil {
        return &gr, nil, err
    }

    return &gr, factors, nil
}

func (r *Repository) SaveGrowthRequest(gr *ds.GrowthRequest) error {
    return r.db.Save(gr).Error
}

func (r *Repository) DeleteGrowthRequest(id uint, creatorID uint) error {
	result := r.db.Model(&ds.GrowthRequest{}).
		Where("id = ? AND creator_id = ? AND status = ?", id, creatorID, "черновик").
		Updates(map[string]interface{}{
			"status":      "удалён",
			"date_finish": time.Now(),
			"date_update": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("growth request not found or cannot be deleted")
	}

	return nil
}
