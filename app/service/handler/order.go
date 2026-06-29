package handler

import (
	"app/pkg/utils"
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model/request"
	"app/service/model/response"

	"github.com/gin-gonic/gin"
)

func orderRouter(g *gin.RouterGroup) {
	OrderGroup := g.Group("/batch_order", middleware.AuthMiddleware())
	{
		// 记账单
		OrderGroup.POST("/temp/create", TempOrderCreate)
		// 下单
		OrderGroup.POST("/create", OrderCreate)
		OrderGroup.POST("/update", OrderUpdate)
		OrderGroup.POST("/update_status", OrderUpdateStatus)
		OrderGroup.POST("/detail", OrderDetail)
		OrderGroup.POST("/list", OrderList)
		OrderGroup.POST("/share", OrderShared)
		OrderGroup.POST("/latest_by_goods", OrderGoodsLatest)
		OrderGroup.POST("/credit/list", CreditList)
		// 生成分享PDF
		OrderGroup.POST("/share/generate", ShareGeneratePDF)
		// 生成分享HTML
		OrderGroup.POST("/share/generate_html", ShareGenerateHTML)
		// 生成分享图片（PNG，适合移动端微信分享）
		OrderGroup.POST("/share/generate_image", ShareGenerateImage)
		// 获取分享数据（JSON）
		OrderGroup.POST("/share/daily", ShareDailyOrder)
	}
}

// 码单
func TempOrderCreate(c *gin.Context) {
	order := logic.NewOrderLogic(c)
	if err := c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := order.TempCreate(); err != nil {
		common.Response(c, err, nil)
		return
	}
	// go order.Record(false, model.HistoryStepOrder, model.PayFeild{})
	common.Response(c, nil, order)
}

func OrderCreate(c *gin.Context) {
	var form request.OrderReq
	var err error
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	form.OwnerUser = common.GetUserUUID(c)
	order := logic.NewOrderLogic(c)
	err = order.Create(form)
	common.Response(c, err, order)
}

func OrderUpdate(c *gin.Context) {
	var err error
	var form request.OrderReq
	if err = c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}
	form.OwnerUser = common.GetUserUUID(c)
	orderLogic := logic.NewOrderLogic(c)

	err = orderLogic.Update(form)
	common.Response(c, err, orderLogic)
}

func OrderUpdateStatus(c *gin.Context) {
	order := logic.NewOrderLogic(c)
	var err error
	if err = c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	err = order.UpdateStatus()
	common.Response(c, err, order)
}

func OrderShared(c *gin.Context) {
	order := logic.NewOrderLogic(c)
	var err error
	if err = c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	err = order.Shared()
	common.Response(c, err, order)
}

func OrderDetail(c *gin.Context) {
	var payLoad request.BatchOrderDetailReq
	var err error
	if err = c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewOrderLogic(c)
	err = order.FromUUID(payLoad.UUID)
	common.Response(c, err, order)
}

func OrderGoodsLatest(c *gin.Context) {

	var body request.BatchOrderGoodsLatest
	var err error
	order := logic.NewOrderLogic(c)
	if err = c.ShouldBind(&body); err != nil {
		common.Response(c, err, nil)
		return
	}

	var goodPriceObjs []response.BatchOrderGoodsOrderRsp
	goodPriceObjs, err = order.FindLatestGoods(body.GoodsUUIDList)
	common.Response(c, err, goodPriceObjs)
}

func OrderList(c *gin.Context) {
	var payload request.BatchOrderListReq
	order := logic.NewOrderLogic(c)
	if err := c.ShouldBind(&payload); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodPriceObjs, err := order.List(payload.UserUUID, payload.StartTime, payload.EndTime, payload.Status, utils.DefaultSetLimitCond(payload.LimitCond))
	common.Response(c, err, goodPriceObjs)
}

func OrderGoodsList(c *gin.Context) {
	var payLoad request.BatchGoodsListReq
	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewOrderLogic(c).GoodsList(payLoad)
	common.Response(c, err, rsp)
}

// CreditList 赊欠列表接口
func CreditList(c *gin.Context) {
	var payLoad request.CreditListReq
	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewOrderLogic(c).CreditList(payLoad)
	common.Response(c, err, rsp)
}

// ShareDailyOrder 查询今日订单分享数据（JSON）
func ShareDailyOrder(c *gin.Context) {
	var req request.ShareDailyOrderReq
	if err := c.ShouldBind(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewOrderLogic(c).ShareDailyOrder(req)
	common.Response(c, err, rsp)
}

// ShareGeneratePDF 生成今日账单PDF，直接返回PDF文件流（无需存档，多台服务器兼容）
func ShareGeneratePDF(c *gin.Context) {
	var req request.ShareDailyOrderReq
	if err := c.ShouldBind(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	pdfBytes, err := logic.NewOrderLogic(c).GenerateSharePDF(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	c.Header("Content-Disposition", "inline; filename=\"receipt.pdf\"")
	c.Data(200, "application/pdf", pdfBytes)
}

// ShareGenerateHTML 生成今日账单HTML，直接返回HTML内容
func ShareGenerateHTML(c *gin.Context) {
	var req request.ShareDailyOrderReq
	if err := c.ShouldBind(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	html, err := logic.NewOrderLogic(c).GenerateShareHTML(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

// ShareGenerateImage 生成今日账单图片（PNG），适合移动端微信分享
func ShareGenerateImage(c *gin.Context) {
	var req request.ShareDailyOrderReq
	if err := c.ShouldBind(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	imageBytes, err := logic.NewOrderLogic(c).GenerateShareImage(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	c.Header("Content-Disposition", "inline; filename=\"receipt.png\"")
	c.Header("Cache-Control", "no-cache")
	c.Data(200, "image/png", imageBytes)
}
