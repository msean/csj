package cache

import (
	"app/pkg/utils"
	"app/service/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type customerCache struct {
	rdb redis.UniversalClient
	db  *gorm.DB
}

func newCustomerCache(rdb redis.UniversalClient, db *gorm.DB) *customerCache {
	return &customerCache{
		rdb: rdb,
		db:  db,
	}
}

// Hash key based on ownerUser
func (cc *customerCache) customerHashKey(ownerUser string) string {
	return fmt.Sprintf("%s:customer:user:%s", CachePrefix, ownerUser)
}

// CustomerFeildSet with Hash
func (cc *customerCache) CustomerFeildSet(uuid, ownerUser string) (model.CustomerFeild, error) {
	ctx := context.Background()
	hashKey := cc.customerHashKey(ownerUser)

	// Try get from hash first
	cached, err := cc.rdb.HGet(ctx, hashKey, uuid).Result()
	if err == nil {
		var field model.CustomerFeild
		if json.Unmarshal([]byte(cached), &field) == nil {
			return field, nil
		}
	}

	// Cache miss, fetch from DB
	var customer model.Customer
	if err := utils.Find(cc.db, &customer, utils.NewWhereCond("owner_user", ownerUser), utils.WhereUIDCond(uuid)); err != nil {
		return model.CustomerFeild{}, err
	}

	field := model.CustomerFeild{
		CustomerName:  customer.Name,
		CustomerPhone: customer.Phone,
	}

	// Store in hash
	data, _ := json.Marshal(field)
	cc.rdb.HSet(ctx, hashKey, uuid, data)
	cc.rdb.Expire(ctx, hashKey, 30*time.Minute)

	return field, nil
}

// BatchCustomerFeildSet with Hash
func (cc *customerCache) BatchCustomerFeildSet(uuidList []string, ownerUser string) (map[string]model.CustomerFeild, error) {
	ctx := context.Background()
	hashKey := cc.customerHashKey(ownerUser)
	result := make(map[string]model.CustomerFeild, len(uuidList))

	// Batch get from hash - HMGET for multiple fields
	cached, err := cc.rdb.HMGet(ctx, hashKey, uuidList...).Result()
	missingUUIDs := make([]string, 0)

	if err == nil {
		for i, val := range cached {
			if val != nil {
				var field model.CustomerFeild
				if json.Unmarshal([]byte(val.(string)), &field) == nil {
					result[uuidList[i]] = field
					continue
				}
			}
			missingUUIDs = append(missingUUIDs, uuidList[i])
		}
	} else {
		missingUUIDs = uuidList
	}

	// Fetch missing items from DB
	if len(missingUUIDs) > 0 {
		var customers []model.Customer
		if err := utils.Find(cc.db, &customers, utils.NewWhereCond("owner_user", ownerUser),
			utils.NewInCondFromString("uid", missingUUIDs)); err != nil {
			return result, err
		}

		// HMSET for batch storing in hash
		hashData := make(map[string]interface{})
		for _, customer := range customers {
			field := model.CustomerFeild{
				CustomerName:  customer.Name,
				CustomerPhone: customer.Phone,
			}
			result[customer.UID] = field
			data, _ := json.Marshal(field)
			hashData[customer.UID] = data
		}

		if len(hashData) > 0 {
			cc.rdb.HMSet(ctx, hashKey, hashData)
			cc.rdb.Expire(ctx, hashKey, 30*time.Minute)
		}
	}

	return result, nil
}

// Cache invalidation with Hash
func (cc *customerCache) InvalidateCustomerCache(uuid, ownerUser string) error {
	ctx := context.Background()
	hashKey := cc.customerHashKey(ownerUser)
	return cc.rdb.HDel(ctx, hashKey, uuid).Err()
}

func (cc *customerCache) InvalidateBatchCustomerCache(uuidList []string, ownerUser string) error {
	ctx := context.Background()
	hashKey := cc.customerHashKey(ownerUser)
	return cc.rdb.HDel(ctx, hashKey, uuidList...).Err()
}

// Clear all cache for a specific user (useful when user deletes account)
func (cc *customerCache) ClearUserCustomerCache(ownerUser string) error {
	ctx := context.Background()
	hashKey := cc.customerHashKey(ownerUser)
	return cc.rdb.Del(ctx, hashKey).Err()
}

// CRUD operations with cache management
func (cc *customerCache) CreateCustomer(customer *model.Customer) error {
	if err := cc.db.Create(customer).Error; err != nil {
		return err
	}
	return cc.InvalidateCustomerCache(customer.UID, customer.OwnerUser)
}

func (cc *customerCache) UpdateCustomer(uuid, ownerUser string, updates map[string]interface{}) error {
	if err := cc.db.Model(&model.Customer{}).Where("uid = ? AND owner_user = ?", uuid, ownerUser).Updates(updates).Error; err != nil {
		return err
	}
	return cc.InvalidateCustomerCache(uuid, ownerUser)
}

func (cc *customerCache) DeleteCustomer(uuid, ownerUser string) error {
	if err := cc.db.Where("uid = ? AND owner_user = ?", uuid, ownerUser).Delete(&model.Customer{}).Error; err != nil {
		return err
	}
	return cc.InvalidateCustomerCache(uuid, ownerUser)
}
