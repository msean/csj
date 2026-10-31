package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	cryptorand "crypto/rand"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	// mysql
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type BaseModel struct {
	UID       string     `gorm:"primary_key;size:64" json:"uuid"`
	CreatedAt time.Time  `json:"createTime"`
	UpdatedAt time.Time  `json:"updateTime"`
	DeletedAt *time.Time `json:"-"`
}

// object must be an pointer type
// func CreateObj(db *gorm.DB, object interface{}) (err error) {
// 	return db.Create(object).Error
// }

// func Find(db *gorm.DB, dst interface{}, conds ...Cond) (err error) {
// 	for _, cond := range conds {
// 		db = cond.Cond(db)
// 	}
// 	err = db.Find(dst).Error
// 	return
// }

// func First(db *gorm.DB, dst interface{}, conds ...Cond) (err error) {
// 	for _, cond := range conds {
// 		db = cond.Cond(db)
// 	}
// 	err = db.First(dst).Error
// 	return
// }

// func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
// 	if m.UID == "" {
// 		id, _ := uuid.NewV4()
// 		m.UID = id.String()
// 	}
// 	return nil
// }

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	// if m.UID == "" {
	// 	t := time.Now()
	// 	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	// 	m.UID = ulid.MustNew(ulid.Timestamp(t), entropy).String()
	// }
	// return nil

	if m.UID != "" {
		return nil
	}

	// 使用 crypto/rand 作为熵源，更安全
	entropy := ulid.Monotonic(cryptorand.Reader, 0)
	t := time.Now()
	m.UID = ulid.MustNew(ulid.Timestamp(t), entropy).String()
	return nil
}

type Array []interface{}

func (a Array) Value() (driver.Value, error) {
	bytes, err := json.Marshal(a)
	return string(bytes), err
}

func (a *Array) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		return json.Unmarshal([]byte(value), a)
	case []byte:
		return json.Unmarshal(value, a)
	default:
		return fmt.Errorf("Array not support")
	}
}

type Map map[string]string

func (m Map) Value() (driver.Value, error) {
	bytes, err := json.Marshal(m)
	return string(bytes), err
}

func (m *Map) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		return json.Unmarshal([]byte(value), m)
	case []byte:
		return json.Unmarshal(value, m)
	default:
		return errors.New("Map not supported")
	}
}

type Ints []int

func (a Ints) Value() (driver.Value, error) {
	bytes, err := json.Marshal(a)
	return string(bytes), err
}

func (a *Ints) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		return json.Unmarshal([]byte(value), a)
	case []byte:
		return json.Unmarshal(value, a)
	default:
		return fmt.Errorf("not support")
	}
}
