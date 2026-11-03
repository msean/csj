package dao

import (
	"time"

	"app/service/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MessageRow keeps database fields separate from the public response.
type MessageRow struct {
	model.MessageCenter
	Source string
	Title  string
}

func visibleOfficial(db *gorm.DB, now time.Time) *gorm.DB {
	return db.Model(&model.OfficialMessage{}).Where("deleted_at IS NULL AND status = 1 AND published_at <= ? AND (expires_at IS NULL OR expires_at > ?)", now, now)
}

func (dao *paymentDao) HasUnreadMessage(db *gorm.DB, owner string) (bool, error) {
	var exists bool
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM message_centers WHERE owner_user = ? AND is_read = 0)", owner).Scan(&exists).Error
	return exists, err
}
func (dao *paymentDao) FindPersonalMessages(db *gorm.DB, owner string, offset, limit int) ([]MessageRow, error) {
	rows := make([]MessageRow, 0)
	err := db.Raw(`SELECT *, 'personal' AS source, '' AS title FROM message_centers WHERE owner_user = ? ORDER BY created_at DESC, uid DESC LIMIT ? OFFSET ?`, owner, limit, offset).Scan(&rows).Error
	return rows, err
}

func (dao *paymentDao) OfficialMessages(db *gorm.DB, now time.Time, offset, limit int) (rows []model.OfficialMessage, total int64, err error) {
	rows = make([]model.OfficialMessage, 0)
	if err = visibleOfficial(db, now).Count(&total).Error; err != nil {
		return
	}
	err = db.Raw(`SELECT o.* FROM official_messages o
 WHERE o.deleted_at IS NULL AND o.status = 1 AND o.published_at <= ? AND (o.expires_at IS NULL OR o.expires_at > ?)
 ORDER BY o.sort DESC, o.published_at DESC, o.uid DESC LIMIT ? OFFSET ?`, now, now, limit, offset).Scan(&rows).Error
	return
}

func (dao *paymentDao) MarkOfficialRead(db *gorm.DB, owner, uid string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var message model.OfficialMessage
		if err := visibleOfficial(tx, time.Now()).Clauses(clause.Locking{Strength: "SHARE"}).Where("uid = ?", uid).First(&message).Error; err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.OfficialMessageRead{UserUUID: owner, MessageUUID: uid, ReadAt: time.Now()}).Error
	})
}

func (dao *paymentDao) MarkAllMessagesRead(db *gorm.DB, owner string) error {
	return db.Model(&model.MessageCenter{}).Where("owner_user = ? AND is_read = 0", owner).Update("is_read", 1).Error
}
