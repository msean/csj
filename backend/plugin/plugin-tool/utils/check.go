package utils

import (
	"fmt"

	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/system"
)

func RegisterApis(apis ...system.SysApi) {
	var count int64
	var apiPaths []string
	for i := range apis {
		apiPaths = append(apiPaths, apis[i].Path)
	}
	global.GVA_MYSQL.Find(&[]system.SysApi{}, "path in (?)", apiPaths).Count(&count)
	if count > 0 {
		return
	}
	err := global.GVA_MYSQL.Create(&apis).Error
	if err != nil {
		fmt.Println(err)
	}
}

func RegisterMenus(menus ...system.SysBaseMenu) {
	var count int64
	var menuNames []string
	parentMenu := menus[0]
	otherMenus := menus[1:]
	for i := range menus {
		menuNames = append(menuNames, menus[i].Name)
	}
	global.GVA_MYSQL.Find(&[]system.SysBaseMenu{}, "name in (?)", menuNames).Count(&count)
	if count > 0 {
		return
	}
	err := global.GVA_MYSQL.Create(&parentMenu).Error
	if err != nil {
		fmt.Println(err)
	}
	for i := range otherMenus {
		pid := parentMenu.ID
		otherMenus[i].ParentId = pid
	}
	err = global.GVA_MYSQL.Create(&otherMenus).Error
	if err != nil {
		fmt.Println(err)
	}
}
