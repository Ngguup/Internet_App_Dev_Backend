package ds

import (
	"database/sql"
	"time"
)

type GrowthRequest struct {
	ID          uint      `gorm:"primaryKey"`
	Status      string    `gorm:"type:varchar(15);not null"`
	DateCreate  time.Time      `gorm:"not null"`
	CreatorID   uint         `gorm:"not null"`

	DateUpdate  time.Time
	DateFinish  sql.NullTime `gorm:"default:null"`
	ModeratorID uint

	CurData     int	     
	StartPeriod time.Time `gorm:"type:date"` 
	EndPeriod   time.Time `gorm:"type:date"` 

	Creator   Users `gorm:"foreignKey:CreatorID"`
	Moderator Users `gorm:"foreignKey:ModeratorID"`
}