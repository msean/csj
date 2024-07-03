package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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
func CreateObj(db *gorm.DB, object interface{}) (err error) {
	return db.Create(object).Error
}

// dst若是primarykey有的话，不需要多余的cond
func Find(db *gorm.DB, dst interface{}, conds ...Cond) (err error) {
	for _, cond := range conds {
		db = cond.Cond(db)
	}
	err = db.Find(dst).Error
	return
}

func First(db *gorm.DB, dst interface{}, conds ...Cond) (err error) {
	for _, cond := range conds {
		db = cond.Cond(db)
	}
	err = db.First(dst).Error
	return
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.UID == "" {
		id, _ := uuid.NewV4()
		m.UID = id.String()
	}
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

type DBConf struct {
	User     string `json:"user"`     // 链接用户名
	Password string `json:"password"` // 链接密码
	Host     string `json:"host"`     // 链接主机
	Port     string `json:"port"`     // 链接端口
	DB       string `json:"db"`       // 链接数据库
	Driver   string `json:"driver"`   // 链接类型
	Charset  string `json:"charset"`  // 字符集
}

func New(conf DBConf) (db *gorm.DB, err error) {

	if conf.Charset == "" {
		conf.Charset = "utf8mb4"
	}

	switch conf.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.DB, conf.Charset)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{ // 表名命名策略
				SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`t_user`
			},
		})
	case "postgres":
		// db, err = gorm.Open(conf.Driver, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		// 	conf.Host, conf.Port, conf.User, conf.DB, conf.Password,
		// ))
	case "sqlite3":
		// db, err = gorm.Open(conf.Driver, conf.DB)
	}

	if err != nil {
		return nil, err
	}

	log.Println("init db")

	return db, err
}
