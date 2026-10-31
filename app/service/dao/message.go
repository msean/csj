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

func unreadOfficial(db *gorm.DB, owner string, now time.Time) *gorm.DB {
	return visibleOfficial(db, now).Where("NOT EXISTS (SELECT 1 FROM official_message_reads r WHERE r.message_uuid = official_messages.uid AND r.user_uuid = ?)", owner)
}

func (dao *paymentDao) CountAllMessages(db *gorm.DB, owner string, now time.Time) (total, unread int64, err error) {
	total, unread, err = dao.CountMessages(db, owner)
	if err != nil {
		return
	}
	var official, officialUnread int64
	if err = visibleOfficial(db, now).Count(&official).Error; err != nil {
		return
	}
	if err = unreadOfficial(db, owner, now).Count(&officialUnread).Error; err != nil {
		return
	}
	return total + official, unread + officialUnread, nil
}

func (dao *paymentDao) HasUnreadMessage(db *gorm.DB, owner string) (bool, error) {
	var exists bool
	if err := db.Raw("SELECT EXISTS(SELECT 1 FROM message_centers WHERE owner_user = ? AND is_read = 0)", owner).Scan(&exists).Error; err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	var ids []string
	err := unreadOfficial(db, owner, time.Now()).Limit(1).Pluck("uid", &ids).Error
	return len(ids) > 0, err
}

func (dao *paymentDao) FindAllMessages(db *gorm.DB, owner string, now time.Time, offset, limit int) ([]MessageRow, error) {
	rows := make([]MessageRow, 0)
	err := db.Raw(`SELECT * FROM (
 SELECT uid, created_at, owner_user, customer_uuid, type, event, content, is_read, related_uuid,
 'personal' AS source, '' AS title, 0 AS pinned, 0 AS sort
 FROM message_centers WHERE owner_user = ?
 UNION ALL
 SELECT o.uid, o.published_at AS created_at, '' AS owner_user, '' AS customer_uuid, 3 AS type,
 '官方消息' AS event, o.content,
 CASE WHEN r.message_uuid IS NULL THEN 0 ELSE 1 END AS is_read, '' AS related_uuid,
 'official' AS source, o.title, 1 AS pinned, o.sort
 FROM official_messages o LEFT JOIN official_message_reads r ON r.message_uuid = o.uid AND r.user_uuid = ?
 WHERE o.deleted_at IS NULL AND o.status = 1 AND o.published_at <= ? AND (o.expires_at IS NULL OR o.expires_at > ?)
 ) messages ORDER BY pinned DESC, sort DESC, created_at DESC, uid DESC LIMIT ? OFFSET ?`, owner, owner, now, now, limit, offset).Scan(&rows).Error
	return rows, err
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
