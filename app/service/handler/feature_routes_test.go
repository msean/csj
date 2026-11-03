package handler

import (
	"app/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewMobileRoutesAreRegisteredAndProtected(t *testing.T) {
	old := global.Global
	global.Global = &global.RunTime{Logger: zap.NewNop()}
	t.Cleanup(func() { global.Global = old })
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/api/csj")
	orderRouter(group)
	messageRouter(group)
	for _, path := range []string{"/api/csj/batch_order/customer_purchases", "/api/csj/message/readAll"} {
		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, path, strings.NewReader("{}")))
		if recorder.Code != 200 || !strings.Contains(recorder.Body.String(), `"code":4001`) {
			t.Fatalf("route missing or not protected %s: %d %s", path, recorder.Code, recorder.Body.String())
		}
	}
}
