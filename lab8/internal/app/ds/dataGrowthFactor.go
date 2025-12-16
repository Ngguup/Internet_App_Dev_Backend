package ds

type DataGrowthFactor struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string  `gorm:"type:varchar(255);not null" json:"title"`
	Image       string  `gorm:"type:varchar(500)" json:"image"`
	Coeff       float64 `gorm:"type:double precision;not null" json:"coeff"`
	Description string  `gorm:"type:varchar(1000)" json:"description"`
	Attribute   string  `gorm:"type:varchar(255)" json:"attribute"`
	IsDelete    bool    `gorm:"type:boolean;default:false" json:"isDelete"`
}

