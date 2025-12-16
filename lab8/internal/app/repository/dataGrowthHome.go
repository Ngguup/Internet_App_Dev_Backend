package repository

import (
	"errors"
	"lab1/internal/app/ds"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) GetGrowthRequestByID_(id int) (ds.GrowthRequest, []ds.DataGrowthFactor, []float64, error) {
	var gr ds.GrowthRequest
	var factors []ds.DataGrowthFactor
	var factorNums []float64

	if err := r.db.
		Where("id = ? AND status = ?", id, "черновик").
		First(&gr).Error; err != nil {
		return ds.GrowthRequest{}, nil, nil, err
	}

	if err := r.db.
		Table("data_growth_factors").
		Joins("JOIN growth_request_data_growth_factors grdf ON grdf.data_growth_factor_id = data_growth_factors.id").
		Where("grdf.growth_request_id = ?", gr.ID).
		Where("data_growth_factors.is_delete = false").
		Order("data_growth_factors.id").
		Find(&factors).Error; err != nil {
		return ds.GrowthRequest{}, nil, nil, err
	}

	for _, f := range factors {
		var factorNum float64
		if err := r.db.
			Table("growth_request_data_growth_factors").
			Select("factor_num").
			Where("growth_request_id = ? AND data_growth_factor_id = ?", gr.ID, f.ID).
			Scan(&factorNum).Error; err != nil {
			return ds.GrowthRequest{}, nil, nil, err
		}
		factorNums = append(factorNums, factorNum)
	}

	return gr, factors, factorNums, nil
}

func (r *Repository) GetAllDataGrowthFactors_() ([]ds.DataGrowthFactor, error) {
	var dataGrowthFactors []ds.DataGrowthFactor
	err := r.db.Where("is_delete = false").Find(&dataGrowthFactors).Error // добавили условие
	if err != nil {
		return nil, err
	}
	return dataGrowthFactors, nil
}

func (r *Repository) GetDataGrowthFactorById_(id int) (*ds.DataGrowthFactor, error) {
	var factor ds.DataGrowthFactor

	err := r.db.
		Where("id = ? AND is_delete = ?", id, false).
		First(&factor).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &factor, nil
}

func (r *Repository) SearchDataGrowthFactorsByName_(name string) ([]ds.DataGrowthFactor, error) {
	var dataGrowthFactors []ds.DataGrowthFactor
	err := r.db.Where("title ILIKE ? and is_delete = ?", "%"+name+"%", false).Find(&dataGrowthFactors).Error 
	if err != nil {
		return nil, err
	}
	return dataGrowthFactors, nil
}

func (r *Repository) GetGrowthRequestData_() (uint, int64) {
	var growthRequestID uint
	var count int64
	creatorID := 1

	err := r.db.Model(&ds.GrowthRequest{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&growthRequestID).Error
	if err != nil {
		return 0, 0
	}

	err = r.db.Model(&ds.GrowthRequestDataGrowthFactor{}).Where("growth_request_id = ?", growthRequestID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_data_growth_factors:", err)
	}

	return growthRequestID, count
}

func (r *Repository) AddDataGrowthFactor_(dataGrowthFactorID uint) error {
	creatorID := 1

	var draft ds.GrowthRequest
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&draft).Error

	if err == gorm.ErrRecordNotFound {
		draft = ds.GrowthRequest{
			Status:      "черновик",
			DateCreate:  time.Now(),
			CreatorID:   uint(creatorID),
		}

		if err := r.db.Create(&draft).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	link := ds.GrowthRequestDataGrowthFactor{
		GrowthRequestID:    draft.ID,
		DataGrowthFactorID: dataGrowthFactorID,
	}

	if err := r.db.Create(&link).Error; err != nil {
		return err
	}

	return nil
}

// func (r *Repository) DeleteDataGrowthFactor_(dataGrowthFactorID uint) error {
// 	err := r.db.Model(&ds.DataGrowthFactor{}).Where("id = ?", dataGrowthFactorID).UpdateColumn("is_delete", true).Error
// 	fmt.Println(dataGrowthFactorID)
// 	if err != nil {
// 		return fmt.Errorf("ошибка при удалении чата с id %d: %w", dataGrowthFactorID, err)
// 	}

// 	return nil
// }
