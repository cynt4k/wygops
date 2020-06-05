package models

import (
	vd "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
)

// Group : Group model
type Group struct {
	gorm.Model
	Name  string  `gorm:"size:255;not null;unique" json:"name"`
	Users []*User `gorm:"man2many:user_groups;" json:"users"`
	Type  string  `gorm:"size:100;not null" json:"type"`
}

// Validate : Validate the group
func (g *Group) Validate() error {
	return vd.ValidateStruct(
		vd.Field(&g.Name, vd.RuneLength(0, 255)),
	)
}
