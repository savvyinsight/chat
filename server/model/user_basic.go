package model

import (
	"errors"
	"regexp"

	"chat/global"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"phone"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64
	HeartbeatTime uint64
	LogoutTime    uint64
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	var users []*UserBasic
	db := global.GVA_DB
	db.Find(&users)
	return users
}

// Validate checks email and phone fields using govalidator.
// It returns an error describing the first invalid field it finds, or nil.
func (u *UserBasic) Validate() error {
	if u.Email != "" {
		if !govalidator.IsEmail(u.Email) {
			return errors.New("invalid email")
		}
	}

	if u.Phone != "" {
		re := regexp.MustCompile(`\D`)
		digits := re.ReplaceAllString(u.Phone, "")
		if digits == "" {
			return errors.New("invalid phone: no digits found")
		}
		if !govalidator.IsNumeric(digits) {
			return errors.New("invalid phone: must contain digits only")
		}
		l := len(digits)
		if l < 7 || l > 15 {
			return errors.New("invalid phone: length must be between 7 and 15 digits")
		}
	}

	return nil
}
