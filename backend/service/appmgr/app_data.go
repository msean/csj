package appmgr

import (
	"fmt"
	"reflect"

	"github.com/msean/csj/backend/global"
	appReq "github.com/msean/csj/backend/model/appmgr/request"
	model "github.com/msean/csj/shared/model"
	"github.com/msean/csj/shared/sharding"
	"gorm.io/gorm"
)

// Resource names and table identifiers are server-owned; requests cannot name tables.
type appResource struct {
	model    any
	table    string
	search   []string
	status   string
	sharded  bool
	customer string
	batch    string
}

var appResources = map[string]appResource{
	"customers":        {model: model.Customer{}, table: "customers", search: []string{"name", "phone", "car_no", "remark"}},
	"goods-categories": {model: model.GoodsCategory{}, table: "goods_categories", search: []string{"name"}},
	"goods":            {model: model.Goods{}, table: "goods", search: []string{"name"}, status: "status"},
	"batches":          {model: model.Batch{}, table: "batches", search: []string{"serial_no"}, status: "status"},
	"orders":           {model: model.BatchOrder{}, table: "batch_orders", search: []string{"uid"}, status: "status", sharded: true, customer: "user_uuid", batch: "batch_uuid"},
	"payments":         {model: model.BatchOrderPay{}, table: "batch_order_pay", search: []string{"uid", "batch_order_uuid", "revoked_reason"}, status: "is_revoked", sharded: true, customer: "customer_uuid"},
	"messages":         {model: model.MessageCenter{}, table: "message_centers", search: []string{"event", "content"}, status: "is_read", customer: "customer_uuid"},
	"shares":           {model: model.ShareRecord{}, table: "share_records", search: []string{"customer_uuid"}, customer: "customer_uuid"},
}

type AppDataService struct{}

func resourceQuery(db *gorm.DB, resource string, q appReq.AppDataSearch) (*gorm.DB, appResource, error) {
	spec, ok := appResources[resource]
	if !ok {
		return nil, spec, fmt.Errorf("不支持的数据类型")
	}
	if spec.sharded && q.OwnerUser == "" {
		return nil, spec, fmt.Errorf("请先选择所属用户")
	}
	table := spec.table
	if spec.sharded {
		table = sharding.GetShardTableName(table, q.OwnerUser)
	}
	query := db.Model(spec.model).Table(table).Where("deleted_at IS NULL")
	if q.OwnerUser != "" {
		query = query.Where("owner_user = ?", q.OwnerUser)
	}
	if q.CustomerUUID != "" {
		if spec.customer == "" {
			return nil, spec, fmt.Errorf("当前数据不支持客户筛选")
		}
		query = query.Where(spec.customer+" = ?", q.CustomerUUID)
	}
	if q.BatchUUID != "" {
		if spec.batch == "" {
			return nil, spec, fmt.Errorf("当前数据不支持批次筛选")
		}
		query = query.Where(spec.batch+" = ?", q.BatchUUID)
	}
	if q.Status != nil {
		if spec.status == "" {
			return nil, spec, fmt.Errorf("当前数据不支持状态筛选")
		}
		query = query.Where(spec.status+" = ?", *q.Status)
	}
	if q.Keyword != "" && len(spec.search) > 0 {
		conditions := db.Session(&gorm.Session{NewDB: true})
		for i, field := range spec.search {
			if i == 0 {
				conditions = conditions.Where(field+" LIKE ?", "%"+q.Keyword+"%")
			} else {
				conditions = conditions.Or(field+" LIKE ?", "%"+q.Keyword+"%")
			}
		}
		query = query.Where(conditions)
	}
	if q.From != nil {
		query = query.Where("created_at >= ?", *q.From)
	}
	if q.To != nil {
		query = query.Where("created_at <= ?", *q.To)
	}
	if q.From != nil && q.To != nil && q.From.After(*q.To) {
		return nil, spec, fmt.Errorf("开始时间不能晚于结束时间")
	}
	return query, spec, nil
}

func (s *AppDataService) List(resource string, q appReq.AppDataSearch) (any, int64, error) {
	return s.list(global.GVA_MYSQL, resource, q)
}
func (s *AppDataService) list(db *gorm.DB, resource string, q appReq.AppDataSearch) (any, int64, error) {
	query, spec, err := resourceQuery(db, resource, q)
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, size := appReq.NormalizePage(q.Page, q.PageSize)
	rows := reflect.New(reflect.SliceOf(reflect.TypeOf(spec.model)))
	rows.Elem().Set(reflect.MakeSlice(rows.Elem().Type(), 0, 0))
	err = query.Order("created_at DESC, uid DESC").Offset((page - 1) * size).Limit(size).Find(rows.Interface()).Error
	return rows.Elem().Interface(), total, err
}

type AppDataDetail struct {
	Record  any            `json:"record"`
	Related map[string]any `json:"related"`
}

func (s *AppDataService) Detail(resource, uid, owner string) (AppDataDetail, error) {
	return s.detail(global.GVA_MYSQL, resource, uid, owner)
}
func (s *AppDataService) detail(db *gorm.DB, resource, uid, owner string) (result AppDataDetail, err error) {
	if uid == "" || owner == "" {
		return result, fmt.Errorf("uid和所属用户不能为空")
	}
	query, spec, err := resourceQuery(db, resource, appReq.AppDataSearch{OwnerUser: owner})
	if err != nil {
		return result, err
	}
	record := reflect.New(reflect.TypeOf(spec.model)).Interface()
	if err = query.Where("uid = ?", uid).First(record).Error; err != nil {
		return result, err
	}
	result.Record = record
	result.Related = make(map[string]any)
	// Each related collection is scoped to the verified record and owner.
	switch resource {
	case "batches":
		var goods []model.BatchGoods
		err = db.Where("batch_uuid = ? AND owner_user = ? AND deleted_at IS NULL", uid, owner).Order("sort_no ASC, uid ASC").Find(&goods).Error
		result.Related["goodsList"] = goods
	case "orders":
		var goods []model.BatchOrderGoods
		if err = db.Table(sharding.GetShardTableName("batch_order_goods", owner)).Where("batch_order_uuid = ? AND owner_user = ? AND deleted_at IS NULL", uid, owner).Order("created_at ASC, uid ASC").Find(&goods).Error; err != nil {
			return result, err
		}
		result.Related["goodsList"] = goods
		var history []model.BatchOrderHistory
		err = db.Where("batch_order_uuid = ? AND deleted_at IS NULL", uid).Order("created_at ASC, uid ASC").Find(&history).Error
		result.Related["history"] = history
	case "payments":
		pay := record.(*model.BatchOrderPay)
		var customer model.Customer
		lookupErr := db.Where("uid = ? AND owner_user = ? AND deleted_at IS NULL", pay.CustomerUUID, owner).First(&customer).Error
		if lookupErr != nil && lookupErr != gorm.ErrRecordNotFound {
			return result, lookupErr
		}
		if lookupErr == nil {
			result.Related["customer"] = customer
		}
	}
	return result, err
}
