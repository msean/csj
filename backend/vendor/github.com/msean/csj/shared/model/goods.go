package model

type GoodsCategory struct {
	BaseModel
	OwnerUser string `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
	Name      string `gorm:"column:name;comment:客户名字" json:"name"`
}

type Goods struct {
	BaseModel
	CategoryID string  `gorm:"column:category_id;comment:所属分类" json:"categoryID"`
	OwnerUser  string  `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
	Name       string  `gorm:"column:name;comment:客户名字" json:"name"`
	Typ        int     `gorm:"column:type;comment:类型 1 定装 2散装" json:"type"`
	Price      float64 `gorm:"column:price;type:decimal(10,2);comment:单价" json:"price"`   // 散装多少钱/斤 定装多少钱一件
	Weight     float64 `gorm:"column:weight;type:decimal(10,2);comment:重量" json:"weight"` // 定装是多少斤
	Status     int     `gorm:"column:status;comment:状态" json:"status"`
}
