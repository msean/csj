package appmgr

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	"gorm.io/gorm"
)

func TestSupportSingletonAndClear(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&appmgr.SupportConfig{}); err != nil {
		t.Fatal(err)
	}
	old := global.GVA_MYSQL
	global.GVA_MYSQL = db
	t.Cleanup(func() { global.GVA_MYSQL = old; sqlDB, _ := db.DB(); _ = sqlDB.Close() })
	api := OfficialMessageApi{}
	save := func(body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/officialMessage/supportConfig", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		api.SaveSupportConfig(c)
		return w
	}
	for _, body := range []string{`{"id":99,"phone":"4001234567","wechat":"support"}`, `{"phone":"","wechat":"","description":""}`} {
		if w := save(body); !strings.Contains(w.Body.String(), `"code":0`) {
			t.Fatal(w.Body.String())
		}
	}
	var rows []appmgr.SupportConfig
	if err = db.Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].ID != 1 || rows[0].Phone != "" || rows[0].Wechat != "" {
		t.Fatalf("singleton/clear failed: %+v", rows)
	}
	w := save(`{"description":"` + strings.Repeat("长", 2001) + `"}`)
	if strings.Contains(w.Body.String(), `"code":0`) {
		t.Fatal("oversized config accepted")
	}
}
