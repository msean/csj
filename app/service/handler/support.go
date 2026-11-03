package handler

import (
	"app/global"
	"app/service/common"
	"app/service/model"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SupportConfig(c *gin.Context) {
	var config model.SupportConfig
	err := global.Global.DB.Where("id = ?", 1).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	common.Response(c, err, gin.H{"phone": config.Phone, "wechat": config.Wechat, "description": config.Description})
}
