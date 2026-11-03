package appmgr

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	"github.com/msean/csj/backend/model/common/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"unicode/utf8"
)

func (api *OfficialMessageApi) GetSupportConfig(c *gin.Context) {
	var config appmgr.SupportConfig
	err := global.GVA_MYSQL.Where("id = ?", 1).First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.FailWithMessage("读取客服配置失败", c)
		return
	}
	response.OkWithData(config, c)
}
func (api *OfficialMessageApi) SaveSupportConfig(c *gin.Context) {
	var config appmgr.SupportConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.FailWithMessage("配置格式错误", c)
		return
	}
	config.ID = 1
	config.Phone = strings.TrimSpace(config.Phone)
	config.Wechat = strings.TrimSpace(config.Wechat)
	config.Description = strings.TrimSpace(config.Description)
	if utf8.RuneCountInString(config.Phone) > 50 || utf8.RuneCountInString(config.Wechat) > 100 || utf8.RuneCountInString(config.Description) > 2000 {
		response.FailWithMessage("客服配置内容过长", c)
		return
	}
	if err := global.GVA_MYSQL.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, DoUpdates: clause.AssignmentColumns([]string{"phone", "wechat", "description"})}).Create(&config).Error; err != nil {
		response.FailWithMessage("保存客服配置失败", c)
		return
	}
	response.OkWithMessage("保存成功", c)
}
