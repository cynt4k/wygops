package models

import (
	"regexp"
	"time"

	vd "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Device : Device model
type Device struct {
	ID          uint      `gorm:"primary_key"`
	UserID      uint      `gorm:"" json:"userId"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	PrivateKey  string    `gorm:"size:255;not null" json:"privateKey"`
	PublicKey   string    `gorm:"size:255;not null" json:"publicKey"`
	IPv4Address string    `gorm:"size:15;not null" json:"ipv4Address"`
	IPv6Address string    `gorm:"size:45;not null" json:"ipv6Address"`
	CreatedAt   time.Time `gorm:"precision:6" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"precision:6" json:"updatedAt"`
}

// TableName : Get the database table name
func (d *Device) TableName() string {
	return "devices"
}

func checkKey(value interface{}) error {
	keyRaw, _ := value.(string)
	_, err := wgtypes.ParseKey(keyRaw)
	if err != nil {
		return err
	}
	return nil
}

// Validate : Validate a device model
func (d *Device) Validate() error {
	return vd.ValidateStruct(
		vd.Field(&d.IPv4Address, vd.Match(regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`))), // nolint:lll
		vd.Field(&d.IPv6Address, vd.Match(regexp.MustCompile(`([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`))), // nolint:lll,staticcheck
		vd.Field(&d.PublicKey, vd.By(checkKey)),
		vd.Field(&d.IPv4Address, is.IPv4),
		vd.Field(&d.IPv6Address, is.IPv6),
	)
}
