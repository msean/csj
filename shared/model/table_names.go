package model

// App tables are independent of the admin database's naming strategy.
func (User) TableName() string                { return "users" }
func (Customer) TableName() string            { return "customers" }
func (GoodsCategory) TableName() string       { return "goods_categories" }
func (Goods) TableName() string               { return "goods" }
func (BatchGoods) TableName() string          { return "batch_goods" }
func (BatchOrder) TableName() string          { return "batch_orders" }
func (BatchOrderGoods) TableName() string     { return "batch_order_goods" }
func (BatchOrderOpr) TableName() string       { return "batch_order_oprs" }
func (BatchOrderPay) TableName() string       { return "batch_order_pays" }
func (BatchOrderHistory) TableName() string   { return "batch_order_histories" }
func (MessageCenter) TableName() string       { return "message_centers" }
func (OfficialMessage) TableName() string     { return "official_messages" }
func (OfficialMessageRead) TableName() string { return "official_message_reads" }
func (ShareRecord) TableName() string         { return "share_records" }
func (Sms) TableName() string                 { return "sms" }
