package ds

import (
	"lab1/internal/app/role"
)

type Users struct {
	ID          uint   `gorm:"primary_key;autoIncrement" json:"id"`
	Login       string `gorm:"type:varchar(25);unique;not null" json:"login"`
	Password    string `gorm:"type:varchar(100);not null" json:"-"`
	// IsModerator bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
	Role role.Role `sql:"type:string;"`

}

// type User struct {
// 	UUID uuid.UUID `gorm:"type:uuid"`
// 	Name string    `json:"name"`
// 	Role role.Role `sql:"type:string;"`
// 	Pass string
// }

