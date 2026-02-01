package service

import (
	"chat/global"
	"chat/model"
	"fmt"
	"time"
)

// SaveMessage persists a message and returns any error.
func SaveMessage(m *model.Message) error {
	if m == nil {
		return fmt.Errorf("nil message")
	}
	return global.GVA_DB.Create(m).Error
}

// AckMessage marks a message as delivered/acked by id.
func AckMessage(messageID uint) error {
	now := time.Now()
	result := global.GVA_DB.Model(&model.Message{}).Where("id = ?", messageID).Updates(map[string]interface{}{
		"delivered":    true,
		"delivered_at": &now,
	})
	return result.Error
}
