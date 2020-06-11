package models

import (
	"html"
	"strings"
	"time"

	"github.com/cynt4k/wygops/pkg/util/cryptutil"
	"github.com/cynt4k/wygops/pkg/util/randutil"
	vd "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

// User : User model
type User struct {
	ID        uint      `gorm:"primary_key"`
	Username  string    `gorm:"size:100;not null;unique" json:"username"`
	Password  string    `gorm:"size:100;not null" json:"password"`
	Cipher    string    `gorm:"size:255;not null" json:"cipher"`
	FirstName string    `gorm:"size:100;not null" json:"firstName"`
	LastName  string    `gorm:"size 100;not null" json:"lastName"`
	Mail      string    `gorm:"size 100;not null" json:"mail"`
	Type      string    `gorm:"size 50;not null" json:"type"`
	Devices   []*Device `gorm:"foreignkey:userId" json:"devices"`
	CreatedAt time.Time `gorm:"precision:6" json:"createdAt"`
	UpdatedAt time.Time `gorm:"precision:6" json:"updatedAt"`
}

// TableName : Get the table name
func (*User) TableName() string {
	return "user"
}

// Hash : Hash a string with an bcrypt hash function
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword : Verify the hashed password with the to be checked
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// EncryptCipher : Encrypt the cipher with the password
func EncryptCipher(cipherText string, password string) (string, error) {
	return cryptutil.EncryptStringToBase64(cipherText, password)
}

// DecryptCipher : Decrypt the cipher with the password
func DecryptCipher(cipherText string, password string) (string, error) {
	return cryptutil.DecryptBase64ToString(cipherText, password)
}

// BeforeSave : Hook before saving the user
func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if u.Cipher == "" {
		u.Cipher = randutil.RandStringRunes(64)
	}
	cipher, err := EncryptCipher(u.Cipher, u.Password)

	if err != nil {
		return err
	}
	if u.Type != "ldap" {
		u.Password = string(hashedPassword)
	}
	u.Cipher = cipher
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
