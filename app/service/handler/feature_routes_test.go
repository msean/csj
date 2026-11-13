package handler

import (
	"app/global"
	"app/service/model/request"
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
	for _, path := range []string{"/api/csj/batch_order/customer_purchases", "/api/csj/batch_order/print/ticket", "/api/csj/message/readAll"} {
		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, path, strings.NewReader("{}")))
		if recorder.Code != 200 || !strings.Contains(recorder.Body.String(), `"code":4001`) {
			t.Fatalf("route missing or not protected %s: %d %s", path, recorder.Code, recorder.Body.String())
		}
	}
}

func TestOrderTicketParametersRequireFeedLinesAndAcceptExplicitZero(t *testing.T) {
	for _, input := range []struct {
		body  string
		valid bool
	}{
		{`{"orderUUID":"order","paperWidth":80,"printMode":"image","feedLines":0}`, true},
		{`{"orderUUID":"order","paperWidth":58,"printMode":"image","feedLines":5}`, true},
		{`{"orderUUID":"order","paperWidth":80,"printMode":"image"}`, false},
		{`{"orderUUID":"order","paperWidth":80,"printMode":"text","feedLines":0}`, false},
		{`{"orderUUID":"order","paperWidth":80,"printMode":"image","feedLines":2}`, false},
		{`{"orderUUID":"order","paperWidth":60,"printMode":"image","feedLines":0}`, false},
		{`{"orderUUID":"order","paperWidth":80,"printMode":"pdf","feedLines":0}`, false},
	} {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(input.body))
		var req request.OrderTicketReq
		if err := ctx.ShouldBindJSON(&req); (err == nil) != input.valid {
			t.Errorf("validation for %s: %v", input.body, err)
		}
	}
}
