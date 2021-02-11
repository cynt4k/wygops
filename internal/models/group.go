package models

import (
	"time"

	vd "github.com/go-ozzo/ozzo-validation/v4"
)

// Group : Group model
type Group struct {
	ID        uint         `gorm:"primary_key"`
	Name      string       `gorm:"size:255;not null;unique" json:"name"`
	Users     []*UserGroup `gorm:"association_autoupdate:false;association_autocreate:false;preload:false;foreignkey:GroupID" json:"users"` // nolint:lll
	Type      string       `gorm:"size:100;not null" json:"type"`
	CreatedAt time.Time    `gorm:"precision:6" json:"createdAt"`
	UpdatedAt time.Time    `gorm:"precision:6" json:"updatedAt"`
}

// TableName : Get the database table name
func (g *Group) TableName() string {
	return "groups"
}

// Validate : Validate the group
func (g *Group) Validate() error {
	return vd.ValidateStruct(
		vd.Field(&g.Name, vd.RuneLength(0, 255)),
	)
}

// IsMember : Check if user is in group
func (g *Group) IsMember(userID uint) bool {
	for _, u := range g.Users {
		if u.UserID == userID {
			return true
		}
	}
	return false
}

// UserGroup : Mapping user to groups
type UserGroup struct {
	GroupID uint `gorm:"not null;primary_key;auto_increment:false" json:"groupId"`
	UserID  uint `gorm:"not null;primary_key;auto_increment:false" json:"userId"`

	Group Group `gorm:"association_autoupdate:false;association_autocreate:false;preload:false;foreignkey:GroupID" json:"groups"` // nolint:lll
	User  User  `gorm:"association_autoupdate:false;association_autocreate:false;preload:false;foreignkey:UserID" json:"users"`   // nolint:lll
}

// TableName : Get the table name
func (*UserGroup) TableName() string {
	return "user_group"
}
