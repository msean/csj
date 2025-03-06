package model

import (
	"app/global"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	// mysql
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type BaseModel struct {
	UID       int64      `gorm:"primaryKey;autoIncrement:false" json:"-"`
	CreatedAt time.Time  `json:"createTime"`
	UpdatedAt time.Time  `json:"updateTime"`
	DeletedAt *time.Time `json:"-"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	stmt := tx.Statement
	if stmt == nil || stmt.Table == "" {
		return fmt.Errorf("cannot determine table name")
	}
	tableName := stmt.Table

	node := GetOrCreateGenerator(tableName, int64(global.Global.Node()))

	if m.UID == 0 {
		m.UID, err = node.GenerateID()
	}

	return nil
}

func (b BaseModel) MarshalJSON() ([]byte, error) {
	type Alias BaseModel
	return json.Marshal(&struct {
		UIDCompatible string `json:"uuid"`
		*Alias
	}{
		UIDCompatible: strconv.FormatInt(b.UID, 10),
		Alias:         (*Alias)(&b),
	})
}
