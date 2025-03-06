package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const (
	CommonCondOwnerUser = "owner_user"
	CommonCondUID       = "uid"
	CommonCondName      = "name"
	CommonCondPhone     = "phone"
	CommonCondSerialNo  = "serial_no"
)

const (
	LikeTypeLeft    = 1
	LikeTypeRight   = 2
	LikeTypeBetween = 3
)

type (
	Cond interface {
		Cond(db *gorm.DB) *gorm.DB
	}
	BaseCond struct {
		cond string
	}
	WhereCond struct {
		key   string
		value any
	}
	WhereLikeCond struct {
		key   string
		value any
	}
	CmpCond struct {
		key    string
		symbol string
		val    any
	}
	OrderCond struct {
		orders []string
	}
	InCond struct {
		key   string
		value []any
	}
	LimitCond struct {
		PageCount int `json:"pageCount"`
		Page      int `json:"page"`
	}
)

func NewBaseCond(cond string) BaseCond {
	return BaseCond{cond}
}

func (base BaseCond) Cond(db *gorm.DB) *gorm.DB {
	db = db.Where(base.cond)
	return db
}

func NewWhereCond(key string, value any) WhereCond {
	return WhereCond{
		key:   key,
		value: value,
	}
}

func (w WhereCond) Cond(db *gorm.DB) *gorm.DB {
	db = db.Where(fmt.Sprintf("%s=?", w.key), w.value)
	return db
}

func NewOrderCond(orders ...string) OrderCond {
	return OrderCond{
		orders: orders,
	}
}

func (o OrderCond) Cond(db *gorm.DB) *gorm.DB {
	for _, order := range o.orders {
		db = db.Order(order)
	}
	return db
}

func DefaultSetLimitCond(in LimitCond) (out LimitCond) {
	out = in
	if out.Page == 0 {
		out.Page = 1
	}
	if out.PageCount == 0 {
		out.PageCount = 10
	}
	return
}

func (l LimitCond) Cond(db *gorm.DB) *gorm.DB {
	if l.Page != 0 && l.PageCount != 0 {
		db = db.Limit(l.PageCount).Offset((l.Page - 1) * l.PageCount)
	}
	return db
}

func (wl WhereLikeCond) Cond(db *gorm.DB) *gorm.DB {
	db = db.Where(fmt.Sprintf("%s like %s", wl.key, wl.value))
	return db
}

func NewInCondFromString(key string, v []string) InCond {
	value := []any{}
	for _, _v := range v {
		value = append(value, _v)
	}

	return InCond{
		key:   key,
		value: value,
	}
}

func (i InCond) Cond(db *gorm.DB) *gorm.DB {
	db = db.Where(fmt.Sprintf("%s in (?)", i.key), i.value)
	return db
}

func WhereOwnerUserCond(ownerUser int64) Cond {
	return NewWhereCond(CommonCondOwnerUser, ownerUser)
}

func WhereUIDCond(UID int64) Cond {
	return NewWhereCond(CommonCondUID, UID)
}

func WhereSerialNoCond(serialNo string) Cond {
	return NewWhereCond(CommonCondSerialNo, serialNo)
}
func WhereNameCond(name string) Cond {
	return NewWhereCond(CommonCondName, name)
}

func WherePhoneCond(phone string) Cond {
	return NewWhereCond(CommonCondPhone, phone)
}

func InUIDCondFromString(UIDList []string) Cond {
	return NewInCondFromString(CommonCondUID, UIDList)
}

func UpdateOrderDescCond() Cond {
	return NewOrderCond("updated_at desc")
}

func CreatedOrderAscCond() Cond {
	return NewOrderCond("created_at desc")
}

func CreatedOrderDescCond() Cond {
	return NewOrderCond("created_at desc")
}

func UpdateOrderAscCond() Cond {
	return NewOrderCond("updated_at")
}

func NewWhereLikeCond(key, value string, likeType int) Cond {
	var likeValue string
	switch likeType {
	case LikeTypeLeft:
		likeValue = fmt.Sprintf("'%s%s'", "%", value)
	case LikeTypeRight:
		likeValue = fmt.Sprintf("'%s%s'", value, "%")
	case LikeTypeBetween:
		likeValue = fmt.Sprintf("'%s%s%s'", "%", value, "%")
	}
	return WhereLikeCond{
		key:   key,
		value: likeValue,
	}
}

func NewOrLikeCond(value string, likeType int, keys ...string) Cond {
	var cond string
	for i, key := range keys {
		var likevalue string
		switch likeType {
		case LikeTypeLeft:
			likevalue = fmt.Sprintf("'%s%s'", "%", value)
		case LikeTypeRight:
			likevalue = fmt.Sprintf("'%s%s'", value, "%")
		case LikeTypeBetween:
			likevalue = fmt.Sprintf("'%s%s%s'", "%", value, "%")
		}
		if i == 0 {
			cond = fmt.Sprintf("%s like %s", key, likevalue)
		} else {
			cond += fmt.Sprintf("or %s like %s", key, likevalue)
		}
	}
	return BaseCond{cond}
}

func NewCmpCond(key, symbol string, val any) Cond {
	return CmpCond{
		key:    key,
		symbol: symbol,
		val:    val,
	}
}

func (cond CmpCond) Cond(db *gorm.DB) *gorm.DB {
	db = db.Where(fmt.Sprintf("%s %s '%v'", cond.key, cond.symbol, cond.val))
	return db
}

// object must be an pointer type
func GormCreateObj(db *gorm.DB, object interface{}) (err error) {
	return db.Create(object).Error
}

func GormFind(db *gorm.DB, dst interface{}, conds ...Cond) (err error) {
	for _, cond := range conds {
		db = cond.Cond(db)
	}
	err = db.Find(dst).Error
	return
}

func GormFirst(db *gorm.DB, dst interface{}, conds ...Cond) (err error) {
	for _, cond := range conds {
		db = cond.Cond(db)
	}
	err = db.First(dst).Error
	return
}

type GormArray []interface{}

func (a GormArray) Value() (driver.Value, error) {
	bytes, err := json.Marshal(a)
	return string(bytes), err
}

func (a *GormArray) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		return json.Unmarshal([]byte(value), a)
	case []byte:
		return json.Unmarshal(value, a)
	default:
		return fmt.Errorf("Array not support")
	}
}

type GormMap map[string]string

func (m GormMap) Value() (driver.Value, error) {
	bytes, err := json.Marshal(m)
	return string(bytes), err
}

func (m *GormMap) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		return json.Unmarshal([]byte(value), m)
	case []byte:
		return json.Unmarshal(value, m)
	default:
		return errors.New("Map not supported")
	}
}

type GormInts []int

func (a GormInts) Value() (driver.Value, error) {
	bytes, err := json.Marshal(a)
	return string(bytes), err
}

func (a *GormInts) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		return json.Unmarshal([]byte(value), a)
	case []byte:
		return json.Unmarshal(value, a)
	default:
		return fmt.Errorf("not support")
	}
}
