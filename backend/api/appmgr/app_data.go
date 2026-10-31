package appmgr

import (
	"github.com/gin-gonic/gin"
	appReq "github.com/msean/csj/backend/model/appmgr/request"
	"github.com/msean/csj/backend/model/common/response"
)

type AppDataApi struct{}

func (api *AppDataApi) List(c *gin.Context) {
	var q appReq.AppDataSearch
	if err := c.ShouldBindQuery(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	q.Page, q.PageSize = appReq.NormalizePage(q.Page, q.PageSize)
	list, total, err := appDataService.List(c.Param("resource"), q)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{List: presentAppRecord(list), Total: total, Page: q.Page, PageSize: q.PageSize}, c)
}
func (api *AppDataApi) Detail(c *gin.Context) {
	detail, err := appDataService.Detail(c.Param("resource"), c.Query("uid"), c.Query("ownerUser"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(presentAppRecord(detail), c)
}
