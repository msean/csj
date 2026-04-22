package utils

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	LikeTypeLeft    = 1
	LikeTypeRight   = 2
	LikeTypeBetween = 3
)

const (
	CommonIDKey = "id"
)

type Cond interface {
	Cond(db *gorm.DB) *gorm.DB
}

type (
	BaseCond struct {
		cond string
	}
	CmpCond struct {
		key    string
		value  any
		symbol string
	}
	WhereCond struct {
		key   string
		value any
	}
	InCond struct {
		key    string
		values []any
	}
	NotInCond struct {
		key    string
		values []any
	}
	OrderByCond struct {
		order string
	}
	WhereLikeCond struct {
		key   string
		value any
	}
	LimitCond struct {
		PageSize  int       `json:"page_size" form:"page_size"`
		PageNum   int       `json:"page" form:"page"`
		StartTime time.Time `json:"start_time" form:"start_time"`
		EndTime   time.Time `json:"end_time" form:"end_time"`
	}
	LimitCondAlternative struct {
		PageSize  int    `json:"page_size" form:"page_size"`
		PageNum   int    `json:"page" form:"page"`
		StartTime string `json:"start_time" form:"start_time"`
		EndTime   string `json:"end_time" form:"end_time"`
	}
	PageLimitCond struct {
		PageSize int `json:"page_size" form:"page_size"`
		PageNum  int `json:"page" form:"page"`
	}
	TimeLimitCond struct {
		StartTime time.Time `json:"start_time" form:"start_time"`
		EndTime   time.Time `json:"end_time" form:"end_time"`
	}
	FindFromArrayCond struct {
		Key   string
		Value any
	}
)

// ----------------------------- Convert -----------------------------
func (alt LimitCondAlternative) ToLimitCond() (limitCond LimitCond, err error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	var st, et time.Time
	if alt.StartTime != "" {
		st, err = time.ParseInLocation("2006-01-02 15:04:05", alt.StartTime, loc)
		if err != nil {
			return
		}
	}
	if alt.EndTime != "" {
		et, err = time.ParseInLocation("2006-01-02 15:04:05", alt.EndTime, loc)
		if err != nil {
			return
		}
	}

	limitCond = LimitCond{
		PageSize:  alt.PageSize,
		PageNum:   alt.PageNum,
		StartTime: st,
		EndTime:   et,
	}
	return
}

// ----------------------------- Constructors -----------------------------

func NewWhereLikeCond(key, value string, likeType int) Cond {
	var likeValue string
	switch likeType {
	case LikeTypeLeft:
		likeValue = fmt.Sprintf("%%%s", value)
	case LikeTypeRight:
		likeValue = fmt.Sprintf("%s%%", value)
	case LikeTypeBetween:
		likeValue = fmt.Sprintf("%%%s%%", value)
	}
	return WhereLikeCond{key: key, value: likeValue}
}

func NewBaseCond(cond string) Cond            { return BaseCond{cond: cond} }
func NewWhereCond(key string, value any) Cond { return WhereCond{key: key, value: value} }
func NewCmpCond(key, symbol string, value any) Cond {
	return CmpCond{key: key, value: value, symbol: symbol}
}
func NewNotInCond(key string, values []any) Cond { return NotInCond{key: key, values: values} }
func NewOrderByCond(order string) Cond           { return OrderByCond{order: order} }
func NewInCond(key string, values []any) Cond    { return InCond{key: key, values: values} }

// ----------------------------- Cond Implementations -----------------------------

func (wl WhereLikeCond) Cond(db *gorm.DB) *gorm.DB {
	return db.Where(fmt.Sprintf("%s LIKE ?", wl.key), wl.value)
}

func (w WhereCond) Cond(db *gorm.DB) *gorm.DB {
	return db.Where(fmt.Sprintf("%s = ?", w.key), w.value)
}

func (c CmpCond) Cond(db *gorm.DB) *gorm.DB {
	return db.Where(fmt.Sprintf("%s %s ?", c.key, c.symbol), c.value)
}

func (base BaseCond) Cond(db *gorm.DB) *gorm.DB {
	return db.Where(base.cond)
}

func (o OrderByCond) Cond(db *gorm.DB) *gorm.DB {
	return db.Order(o.order)
}

func (o NotInCond) Cond(db *gorm.DB) *gorm.DB {
	return db.Where(fmt.Sprintf("%s NOT IN ?", o.key), o.values)
}

func (o InCond) Cond(db *gorm.DB) *gorm.DB {
	return db.Where(fmt.Sprintf("%s NOT IN ?", o.key), o.values)
}

func (o FindFromArrayCond) Cond(db *gorm.DB) *gorm.DB {
	cond := fmt.Sprintf("%s LIKE ?", o.Key)
	return db.Where(cond, "%"+fmt.Sprintf("[%v,", o.Value)+"%").
		Or(cond, "%"+fmt.Sprintf(",%v]", o.Value)+"%").
		Or(cond, "%"+fmt.Sprintf(",%v,", o.Value)+"%").
		Or(cond, fmt.Sprintf("[%v]", o.Value))
}

func (l PageLimitCond) Cond(db *gorm.DB) *gorm.DB {
	if l.PageSize > 0 {
		db = db.Limit(l.PageSize).Offset(l.PageSize * max(l.PageNum-1, 0))
	}
	return db
}

func (tL TimeLimitCond) Cond(db *gorm.DB) *gorm.DB {
	if !tL.StartTime.IsZero() {
		db = db.Where("created_at >= ?", tL.StartTime)
	}
	if !tL.EndTime.IsZero() {
		db = db.Where("created_at <= ?", tL.EndTime)
	}
	return db
}

func (l LimitCond) Cond(db *gorm.DB) *gorm.DB {
	if l.PageSize > 0 {
		db = db.Limit(l.PageSize).Offset(l.PageSize * max(l.PageNum-1, 0))
	}
	if !l.StartTime.IsZero() {
		db = db.Where("created_at >= ?", l.StartTime)
	}
	if !l.EndTime.IsZero() {
		db = db.Where("created_at <= ?", l.EndTime)
	}
	return db
}

// ----------------------------- CRUD Helpers -----------------------------

func Find(session *gorm.DB, dst any, conds ...Cond) error {
	conds = append(conds, NewBaseCond("deleted_at is NULL"))
	for _, cond := range conds {
		session = cond.Cond(session)
	}
	return session.Find(dst).Error
}

func Get(session *gorm.DB, dst any, conds ...Cond) (has bool, err error) {
	conds = append(conds, NewBaseCond("deleted_at is NULL"))
	for _, cond := range conds {
		session = cond.Cond(session)
	}
	err = session.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}

func Create(session *gorm.DB, objs ...any) (count int64, err error) {
	tx := session.CreateInBatches(objs, 100)
	return tx.RowsAffected, tx.Error
}

func TotalByConds(session *gorm.DB, model any, conds ...Cond) (total int64, err error) {
	conds = append(conds, NewBaseCond("deleted_at is NULL"))
	for _, cond := range conds {
		session = cond.Cond(session)
	}
	err = session.Model(model).Count(&total).Error
	return
}

func Update(session *gorm.DB, updates map[string]any, conds ...Cond) (rows int64, err error) {
	for _, cond := range conds {
		session = cond.Cond(session)
	}
	tx := session.Updates(updates)
	return tx.RowsAffected, tx.Error
}

// ----------------------------- Utilities -----------------------------

func IDCond(val any) Cond { return NewWhereCond(CommonIDKey, val) }

func MutateLimitCond(in LimitCond) (timeLimit TimeLimitCond, pageLimit PageLimitCond) {
	return TimeLimitCond{StartTime: in.StartTime, EndTime: in.EndTime},
		PageLimitCond{PageSize: in.PageSize, PageNum: max(in.PageNum-1, 0)}
}

func (l PageLimitCond) Mutate() PageLimitCond {
	if l.PageNum >= 1 {
		l.PageNum--
	}
	return l
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func SoftDeleteCond() Cond {
	return NewBaseCond("deleted_at IS NULL")
}
