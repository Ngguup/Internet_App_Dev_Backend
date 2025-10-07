package ds

type DataGrowthFactor struct {
	ID          int     `gorm:"primaryKey;autoIncrement"`                         
	Title       string  `gorm:"type:varchar(255);not null"`        
	Image       string  `gorm:"type:varchar(500); not null"`                  
	Coeff       float64 `gorm:"type:double precision; not null"`              
	Description string  `gorm:"type:varchar(1000)"`                
	IsDelete    bool   `gorm:"type:boolean;default:false"`
}
