package service

import (
	"chat/global"
	"chat/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func GetUserList() (users []model.UserBasic) {
	db := global.GVA_DB.Model(model.UserBasic{})
	db.Find(&users)
	return
}

func CreateUser(user *model.UserBasic) (err error) {
	if err = user.Validate(); err != nil {
		return err
	}

	err = global.GVA_DB.Create(user).Error
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

	updateResult := global.GVA_DB.Model(&model.UserBasic{}).
		Where("id = ?", userID).
		Updates(updateFields)

	if updateResult.Error != nil {
		return updateResult.Error
	}

	return nil
}
