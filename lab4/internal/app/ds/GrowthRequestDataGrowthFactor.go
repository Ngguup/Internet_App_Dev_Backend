package ds

type GrowthRequestDataGrowthFactor struct {
	ID                 uint `gorm:"primaryKey;autoIncrement"`
	GrowthRequestID    uint `gorm:"not null;uniqueIndex:idx_growthRequest_dataGrowthFactor"`
	DataGrowthFactorID uint `gorm:"not null;uniqueIndex:idx_growthRequest_dataGrowthFactor"`

	FactorNum       float64 `gorm:"type:double precision;not null;default:1.0"`

	// связи
	GrowthRequest    GrowthRequest    `gorm:"foreignKey:GrowthRequestID"`
	DataGrowthFactor DataGrowthFactor `gorm:"foreignKey:DataGrowthFactorID"`
}
