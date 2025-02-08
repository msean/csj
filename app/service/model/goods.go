package model

const (
	GoodsTypePack   = 1 // 定装
	GoodsTypeWeight = 2 // 散装
)

type GoodsCategory struct {
	BaseModel
	OwnerUser string  `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
	Name      string  `gorm:"column:name;comment:客户名字" json:"name"`
	Goods     []Goods `gorm:"-" json:"goodsList"`
}

type Goods struct {
	BaseModel
	CategoryID string  `gorm:"column:category_id;comment:所属分类" json:"categoryID"`
	OwnerUser  string  `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
	Name       string  `gorm:"column:name;comment:客户名字" json:"name"`
	Typ        int32   `gorm:"column:type;comment:类型 1 定装 2散装" json:"type"`
	Price      float32 `gorm:"column:price;type:decimal(10,2);comment:单价" json:"price"`
	Weight     float32 `gorm:"column:weight;type:decimal(10,2);comment:重量" json:"weight"`
	Status     int     `gorm:"column:status;comment:状态" json:"status"`
}
