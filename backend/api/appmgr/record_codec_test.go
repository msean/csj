package appmgr

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	model "github.com/msean/csj/shared/model"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAdminRecordCompatibility(t *testing.T) {
	now := time.Now()
	user := model.User{BaseModel: model.BaseModel{UID: "u", CreatedAt: now}, ExperienceDeadline: now, MarketName: "市场"}
	data := presentAppRecord(user).(map[string]any)
	if data["uid"] != "u" || data["createdAt"] == nil || data["experienceDeadline"] == nil || data["marketName"] != "市场" {
		t.Fatalf("admin contract broken: %+v", data)
	}
	if _, ok := data["uuid"]; ok {
		t.Fatal("unexpected app field")
	}
	raw, _ := json.Marshal(user)
	if !strings.Contains(string(raw), `"uuid":"u"`) || !strings.Contains(string(raw), `"vipRemainDays"`) {
		t.Fatal("app contract changed")
	}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("PUT", "/", strings.NewReader(`{"uid":"existing","title":"公告","createTime":"2000-01-01T00:00:00Z"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	var message model.OfficialMessage
	if err := bindAppRecord(c, &message); err != nil {
		t.Fatal(err)
	}
	if message.UID != "existing" || message.Title != "公告" || !message.CreatedAt.IsZero() {
		t.Fatalf("input mapping failed: %+v", message)
	}
}
