package dao

import (
	"app/service/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// shareDao 客户分享记录DAO
type shareDao struct{}

func newShareDao() *shareDao {
	return &shareDao{}
}

// UpsertShare 写入/更新客户最近一次分享时间
// 依赖 (owner_user, customer_uuid) 唯一索引，冲突时仅更新 share_time
func (dao *shareDao) UpsertShare(db *gorm.DB, record *model.ShareRecord) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "owner_user"}, {Name: "customer_uuid"}},
		DoUpdates: clause.AssignmentColumns([]string{"share_time", "updated_at"}),
	}).Create(record).Error
}

// TodayRange 返回今天 [00:00:00, 次日00:00:00) 的 unix 时间戳范围
func TodayRange(now time.Time) (startUnix, endUnix int64) {
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return start.Unix(), start.AddDate(0, 0, 1).Unix()
}

// SharedTodayCustomers 查询当天已分享的客户集合，返回 map[customerUUID]true
func (dao *shareDao) SharedTodayCustomers(db *gorm.DB, ownerUser string) (map[string]bool, error) {
	startUnix, endUnix := TodayRange(time.Now())

	var customerUUIDs []string
	err := db.Model(&model.ShareRecord{}).
		Where("owner_user = ? AND share_time >= ? AND share_time < ?", ownerUser, startUnix, endUnix).
		Pluck("customer_uuid", &customerUUIDs).Error
	if err != nil {
		return nil, err
	}

	set := make(map[string]bool, len(customerUUIDs))
	for _, uuid := range customerUUIDs {
		set[uuid] = true
	}
	return set, nil
}
