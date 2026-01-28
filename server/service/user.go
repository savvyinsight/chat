package service

import (
	"chat/global"
	"chat/model"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

func GetUserList() (users []model.UserBasic) {
	db := global.GVA_DB.Model(model.UserBasic{})
	db.Find(&users)
	return
}

func CreateUser(user *model.UserBasic) (err error) {
	db := global.GVA_DB

	var existing model.UserBasic
	if user.Email != "" {
		if err := db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
			return fmt.Errorf("email already registered")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if user.Phone != "" {
		if err := db.Where("phone = ?", user.Phone).First(&existing).Error; err == nil {
			return fmt.Errorf("phone already registered")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	// hash password before storing
	if user.PassWord != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.PassWord), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PassWord = string(hash)
	}

	if err = user.Validate(); err != nil {
		return err
	}

	err = db.Create(user).Error
	return
}

func DeleteUser(user *model.UserBasic) (err error) {
	result := global.GVA_DB.Delete(user)
	if result.Error != nil {
		return result.Error
	}

	// Check actual delete row counts
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func UpdateUser(user *model.UserBasic) error {
	if err := user.Validate(); err != nil {
		return err
	}

	// hash password if provided
	if user.PassWord != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.PassWord), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PassWord = string(hash)
	}

	// check duplicates for email/phone
	var conflict model.UserBasic
	if user.Email != "" {
		if err := global.GVA_DB.Where("email = ? AND id <> ?", user.Email, user.ID).First(&conflict).Error; err == nil {
			return fmt.Errorf("email already registered")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if user.Phone != "" {
		if err := global.GVA_DB.Where("phone = ? AND id <> ?", user.Phone, user.ID).First(&conflict).Error; err == nil {
			return fmt.Errorf("phone already registered")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	result := global.GVA_DB.Model(&model.UserBasic{}).
		Where("id = ?", user.ID).
		Updates(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	return nil
}

func UpdateUserPartial(userID uint, updateFields map[string]interface{}) error {
	// 1. 先检查用户是否存在
	var existingUser model.UserBasic
	result := global.GVA_DB.First(&existingUser, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return result.Error
	}

	// 2. 执行更新
	// Handle password hashing if present in updateFields
	// Accept common JSON keys: "password", "PassWord", and normalize to DB column "pass_word"
	if pwVal, ok := updateFields["password"]; ok {
		if pwStr, ok2 := pwVal.(string); ok2 {
			if pwStr != "" {
				hash, err := bcrypt.GenerateFromPassword([]byte(pwStr), bcrypt.DefaultCost)
				if err != nil {
					return err
				}
				updateFields["pass_word"] = string(hash)
			}
		} else {
			return fmt.Errorf("invalid password format")
		}
		delete(updateFields, "password")
	}
	if pwVal, ok := updateFields["PassWord"]; ok {
		if pwStr, ok2 := pwVal.(string); ok2 {
			if pwStr != "" {
				hash, err := bcrypt.GenerateFromPassword([]byte(pwStr), bcrypt.DefaultCost)
				if err != nil {
					return err
				}
				updateFields["pass_word"] = string(hash)
			}
		} else {
			return fmt.Errorf("invalid password format")
		}
		delete(updateFields, "PassWord")
	}

	// Validate email/phone if present in updateFields
	if emailVal, ok := updateFields["email"]; ok {
		if emailStr, ok2 := emailVal.(string); ok2 {
			temp := existingUser
			temp.Email = emailStr
			if err := temp.Validate(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid email format")
		}
	}
	if phoneVal, ok := updateFields["phone"]; ok {
		if phoneStr, ok2 := phoneVal.(string); ok2 {
			temp := existingUser
			temp.Phone = phoneStr
			if err := temp.Validate(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid phone format")
		}
	}

	// check for duplicates when updating fields
	if emailVal, ok := updateFields["email"]; ok {
		if emailStr, ok2 := emailVal.(string); ok2 {
			var conflict model.UserBasic
			if err := global.GVA_DB.Where("email = ? AND id <> ?", emailStr, userID).First(&conflict).Error; err == nil {
				return fmt.Errorf("email already registered")
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
	}
	if phoneVal, ok := updateFields["phone"]; ok {
		if phoneStr, ok2 := phoneVal.(string); ok2 {
			var conflict model.UserBasic
			if err := global.GVA_DB.Where("phone = ? AND id <> ?", phoneStr, userID).First(&conflict).Error; err == nil {
				return fmt.Errorf("phone already registered")
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
	}

	updateResult := global.GVA_DB.Model(&model.UserBasic{}).
		Where("id = ?", userID).
		Updates(updateFields)

	if updateResult.Error != nil {
		return updateResult.Error
	}

	return nil
}
