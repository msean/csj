package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/model"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommonLogic struct {
	context *gin.Context
	runtime *global.RunTime
}

func NewCommonLogic(context *gin.Context) *CommonLogic {
	return &CommonLogic{
		context: context,
		runtime: global.Global,
	}
}

func (logic *CommonLogic) Init() (result any, err error) {
	result = map[string]interface{}{
		"needUpdate":    false,
		"latestVersion": "",
		"downloadUrl":   "",
		"isForce":       false,
		"description":   "",
	}
	// 1. 获取请求头中的版本号
	currentVersion := logic.context.GetHeader("version")
	if currentVersion == "" {
		return result, errors.New("版本号不能为空")
	}

	// 2. 查询最新的版本记录
	var latestVersion model.AppVersion
	err = logic.runtime.DB.Order("created_at desc").First(&latestVersion).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, nil
		}
		return result, err
	}

	// 3. 比较版本号
	needUpdate := false
	compareResult := utils.CompareVersion(currentVersion, latestVersion.Version)
	if compareResult < 0 {
		needUpdate = true
	}

	// 4. 构造返回结果
	result = map[string]interface{}{
		"needUpdate":    needUpdate,
		"latestVersion": latestVersion.Version,
		"downloadUrl":   latestVersion.DownloadAddr,
		"isForce":       latestVersion.IsForce,
		"description":   latestVersion.Description,
	}

	return result, nil
}
