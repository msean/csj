package appmgr

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
	"gorm.io/gorm"
)

func TestOfficialMessageLifecycle(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&appmgr.OfficialMessage{}); err != nil {
		t.Fatal(err)
	}
	oldDB := global.GVA_MYSQL
	global.GVA_MYSQL = db
	t.Cleanup(func() {
		global.GVA_MYSQL = oldDB
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})
	service := OfficialMessageService{}
	expires := time.Now().Add(time.Hour)
	m := appmgr.OfficialMessage{Title: "公告", Content: "内容", Status: 1, Sort: 10, ExpiresAt: &expires}
	if err = service.CreateOfficialMessage(&m); err != nil {
		t.Fatal(err)
	}
	if m.UID == "" || m.PublishedAt.IsZero() {
		t.Fatal("missing generated UID/publication time")
	}
	// Updating must persist zero values and a cleared expiry, so admins can withdraw messages.
	m.Status, m.Sort, m.ExpiresAt, m.Level = 0, 0, nil, 2
	if err = service.UpdateOfficialMessage(m); err != nil {
		t.Fatal(err)
	}
	got, err := service.GetOfficialMessage(m.UID)
	if err != nil || got.Status != 0 || got.Sort != 0 || got.ExpiresAt != nil || got.Level != 2 {
		t.Fatalf("draft update not persisted: %+v, %v", got, err)
	}
	published := 1
	list, total, err := service.GetOfficialMessageInfoList(appmgrReq.OfficialMessageSearch{Status: &published})
	if err != nil || total != 0 || len(list) != 0 {
		t.Fatalf("draft visible in published filter: %v %d %v", list, total, err)
	}
	draft := 0
	_, total, err = service.GetOfficialMessageInfoList(appmgrReq.OfficialMessageSearch{Status: &draft})
	if err != nil || total != 1 {
		t.Fatalf("draft filter failed: %d %v", total, err)
	}
	if err = service.DeleteOfficialMessage(m.UID); err != nil {
		t.Fatal(err)
	}
	_, total, err = service.GetOfficialMessageInfoList(appmgrReq.OfficialMessageSearch{})
	if err != nil || total != 0 {
		t.Fatalf("deleted message still listed: %d %v", total, err)
	}
}

func TestOfficialMessageValidation(t *testing.T) {
	now := time.Now()
	for _, tc := range []struct {
		name string
		m    appmgr.OfficialMessage
	}{
		{"blank title", appmgr.OfficialMessage{Title: " ", Content: "内容"}},
		{"blank content", appmgr.OfficialMessage{Title: "标题", Content: " "}},
		{"invalid level", appmgr.OfficialMessage{Title: "标题", Content: "内容", Level: 3}},
		{"invalid status", appmgr.OfficialMessage{Title: "标题", Content: "内容", Status: 2}},
		{"expired at publication", appmgr.OfficialMessage{Title: "标题", Content: "内容", PublishedAt: now, ExpiresAt: &now}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if validateOfficialMessage(&tc.m) == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
