package model

// SupportConfig is a singleton managed by the admin, shared by all app users.
type SupportConfig struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Phone       string `json:"phone" gorm:"size:50"`
	Wechat      string `json:"wechat" gorm:"size:100"`
	Description string `json:"description" gorm:"size:2000"`
}
