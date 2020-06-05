package models

import (
	"html"
	"strings"
	"time"

	vd "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User : User model
type User struct {
	gorm.Model
	Username  string   `gorm:"size:100;not null;unique" json:"username"`
	Password  string   `gorm:"size:100;not null" json:"password"`
	FirstName string   `gorm:"size:100;not null" json:"firstName"`
	LastName  string   `gorm:"size 100;not null" json:"lastName"`
	Mail      string   `gorm:"size 100;not null" json:"mail"`
	Type      string   `gorm:"size 50;not null" json:"type"`
	Devices   []Device `gorm:"foreignkey:userId" json:"devices"`
	Groups    []*Group `gorm:"many2many:groups;" json:"user_groups"`
}

// Hash : Hash a string with an bcrypt hash function
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword : Verify the hashed password with the to be checked
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave : Hook before saving the user
func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// Prepare : Prepare the user
func (u *User) Prepare() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Mail = html.EscapeString(strings.TrimSpace(u.Mail))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Validate : Validate if the user is valid
func (u *User) Validate() error {
	return vd.ValidateStruct(
		vd.Field(&u.Username, vd.RuneLength(0, 100)),
		vd.Field(&u.Password, vd.RuneLength(0, 100)),
	)
}
