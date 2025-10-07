package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type DataGrowthFactor struct {
	ID          int
	Title       string
	Image       string
	Coeff       float32
	Description string
}

type GrowthRequest struct {
	ID          uint      
	Status      bool  
	CurData     int	     
	StartPeriod string
	EndPeriod   string
	Components []DataGrowthFactor
}

type GrowthRequestDataGrowthFactor struct {
	ID                 uint 
	GrowthRequestID    uint 
	DataGrowthFactorID uint 

	FactorNum       float64 
}

func (r *Repository) GetGrowthRequestByID(id int) (GrowthRequest, []float64, error) {
	dataGrowthFactors, err := r.GetDataGrowthFactors()
	if err != nil {
		return GrowthRequest{}, nil, err
	}

	grdf := []GrowthRequestDataGrowthFactor{
		{ID: 1, GrowthRequestID: 1, DataGrowthFactorID: uint(dataGrowthFactors[0].ID), FactorNum: 1.0},
		{ID: 2, GrowthRequestID: 1, DataGrowthFactorID: uint(dataGrowthFactors[1].ID), FactorNum: 1.2},

		{ID: 3, GrowthRequestID: 2, DataGrowthFactorID: uint(dataGrowthFactors[3].ID), FactorNum: 0.8},
		{ID: 4, GrowthRequestID: 2, DataGrowthFactorID: uint(dataGrowthFactors[2].ID), FactorNum: 1.5},
		{ID: 5, GrowthRequestID: 2, DataGrowthFactorID: uint(dataGrowthFactors[1].ID), FactorNum: 1.1},
		{ID: 6, GrowthRequestID: 2, DataGrowthFactorID: uint(dataGrowthFactors[0].ID), FactorNum: 1.0},
	}

	growthRequests := []GrowthRequest{
		{
			ID:          1,
			Status:      false,
			CurData:     100,
			StartPeriod: "01.01.24",
			EndPeriod:   "31.12.24",
		},
		{
			ID:          2,
			Status:      true,
			CurData:     200,
			StartPeriod: "01.01.23",
			EndPeriod:   "31.12.23",
		},
	}

	var allFactorNums [][]float64

	for i := range growthRequests {
		var components []DataGrowthFactor
		var factorNums []float64

		for _, rel := range grdf {
			if rel.GrowthRequestID == growthRequests[i].ID {
				for _, f := range dataGrowthFactors {
					if uint(f.ID) == rel.DataGrowthFactorID {
						components = append(components, f)
						factorNums = append(factorNums, rel.FactorNum)
					}
				}
			}
		}

		growthRequests[i].Components = components
		allFactorNums = append(allFactorNums, factorNums)
	}

	return growthRequests[id], allFactorNums[id], nil
}


func (r *Repository) GetDataGrowthFactors() ([]DataGrowthFactor, error) {
	dataGrowthFactors := []DataGrowthFactor{
		{
			ID:    1,
			Title: "Количество пользователей",
			Image: "http://127.0.0.1:9000/pictures/Users.jpg",
			Coeff: 10,
			Description: `
		Показатель, отражающий общее число активных участников системы.
		Он определяет масштаб нагрузки и напрямую влияет на объем генерируемых данных.
		`,
		},
		{
			ID:    2,
			Title: "Частота операций",
			Image: "http://127.0.0.1:9000/pictures/Frequency.jpg",
			Coeff: 20,
			Description: `
	  	Характеристика, показывающая, как часто выполняются вычислительные действия при обработке данных. 
		Она отражает интенсивность загрузки вычислительной системы.
	  `,
		},
		{
			ID:    3,
			Title: "Размер записей",
			Image: "http://127.0.0.1:9000/pictures/Entry.jpg",
			Coeff: 30,
			Description: `
		Характеристика, задающая средний объём одной записи в таблице.
		Чем больше размер, тем быстрее растёт общий объём данных при их накоплении.
	  `,
		},
		{
			ID:    4,
			Title: "Количество файлов",
			Image: "http://127.0.0.1:9000/pictures/Files.jpg",
			Coeff: 40,
			Description: `
		Показатель, определяющий число прикреплённых или загружаемых файлов.
		Он отражает интенсивность использования файлового хранилища в системе.
	  `,
		},
		{
			ID:    5,
			Title: "Количество версий записей",
			Image: "http://127.0.0.1:9000/pictures/Entries.jpg",
			Coeff: 50,
			Description: `
		Характеристика, показывающая, сколько изменений или редакций сохраняется для каждой записи.
		Она определяет скорость роста данных при ведении истории изменений.
	  `,
		},
	}

	if len(dataGrowthFactors) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return dataGrowthFactors, nil
}

func (r *Repository) GetDataGrowthFactor(id int) (DataGrowthFactor, error) {
	dataGrowthFactors, err := r.GetDataGrowthFactors()
	if err != nil {
		return DataGrowthFactor{}, err
	}

	for _, dataGrowthFactor := range dataGrowthFactors {
		if dataGrowthFactor.ID == id {
			return dataGrowthFactor, nil
		}
	}
	return DataGrowthFactor{}, fmt.Errorf("заказ не найден")
}

func (r *Repository) GetDataGrowthFactorsByTitle(title string) ([]DataGrowthFactor, error) {
	dataGrowthFactors, err := r.GetDataGrowthFactors()
	if err != nil {
		return []DataGrowthFactor{}, err
	}

	var result []DataGrowthFactor
	for _, dataGrowthFactor := range dataGrowthFactors {
		if strings.Contains(strings.ToLower(dataGrowthFactor.Title), strings.ToLower(title)) {
			result = append(result, dataGrowthFactor)
		}
	}

	return result, nil
}


