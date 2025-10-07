package repository

import (
	"fmt"
	"lab1/internal/app/ds"
	"gorm.io/gorm"
)

func (r *Repository) DeleteFromDraft(dataGrowthFactorID uint, creatorID uint) error {
	var grds ds.GrowthRequestDataGrowthFactor

	result := r.db.Joins("JOIN growth_requests gr ON gr.id = growth_request_data_growth_factors.growth_request_id").
		Where("gr.creator_id = ? AND gr.status = ? AND growth_request_data_growth_factors.data_growth_factor_id = ?", creatorID, "черновик", dataGrowthFactorID).
		First(&grds)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("entry not found")
		}
		return result.Error
	}

	if err := r.db.Delete(&grds).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateFactorNum(dataGrowthFactorID uint, newFactorNum float64, creatorID uint) error {
	var grds ds.GrowthRequestDataGrowthFactor

	result := r.db.Joins("JOIN growth_requests gr ON gr.id = growth_request_data_growth_factors.growth_request_id").
		Where("gr.creator_id = ? AND gr.status = ? AND growth_request_data_growth_factors.data_growth_factor_id = ?", creatorID, "черновик", dataGrowthFactorID).
		First(&grds)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("entry not found")
		}
		return result.Error
	}

	if err := r.db.Model(&grds).Update("factor_num", newFactorNum).Error; err != nil {
		return err
	}

	return nil
}
