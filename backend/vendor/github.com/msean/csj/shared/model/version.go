package model

type AppVersion struct {
	BaseModel
	Version      string `gorm:"column:version;size:100;not null;comment:版本号" json:"version"`
	DownloadAddr string `gorm:"column:download_addr;size:100;not null;comment:下载地址" json:"downloadAddr"`
	IsForce      bool   `gorm:"column:is_force;default:false;comment:是否强制更新" json:"isForce"`
	Description  string `gorm:"column:description;size:500;comment:更新说明" json:"description"`
}

func (AppVersion) TableName() string {
	return "app_versions"
}
