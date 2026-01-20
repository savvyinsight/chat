package service

import (
	"chat/global"
	"chat/model"
	"fmt"
)

func GetUserList() (users []model.UserBasic) {
	db := global.GVA_DB.Model(model.UserBasic{})
	db.Find(&users)
	return
}

func CreateUser(user *model.UserBasic) (err error) {
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
