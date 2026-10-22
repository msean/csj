package logic

import (
	"app/global"
	"app/service/model"

	"github.com/gin-gonic/gin"
)

type MemberLogic struct {
	context *gin.Context
	runtime *global.RunTime
}

func NewMemberLogic(context *gin.Context) *MemberLogic {
	return &MemberLogic{
		context: context,
		runtime: global.Global,
	}
}

// GetMemberConfigs 获取当前会员收费配置
func (logic *MemberLogic) GetMemberConfigs() (configs []model.MemberConfig, err error) {
	configs = make([]model.MemberConfig, 0)
	err = logic.runtime.DB.Model(&model.MemberConfig{}).Order("type asc").Find(&configs).Error
	return
}
