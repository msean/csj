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

type goodsCache struct {
	rdb redis.UniversalClient
	db  *gorm.DB
}

func newGoodsCache(rdb redis.UniversalClient, db *gorm.DB) *goodsCache {
	return &goodsCache{
		rdb: rdb,
		db:  db,
	}
}

// Hash key based on ownerUser
func (gc *goodsCache) goodsHashKey(ownerUser string) string {
	return fmt.Sprintf("%s:goods:customer:%s", CachePrefix, ownerUser)
}

// model.GoodsFeildSet with Hash
func (gc *goodsCache) GoodsFeildSet(uuid, ownerUser string) (model.GoodsFeild, error) {
	ctx := context.Background()
	hashKey := gc.goodsHashKey(ownerUser)

	// Try get from hash first
	cached, err := gc.rdb.HGet(ctx, hashKey, uuid).Result()
	if err == nil {
		var field model.GoodsFeild
		if json.Unmarshal([]byte(cached), &field) == nil {
			return field, nil
		}
	}

	// Cache miss, fetch from DB
	var _goods model.Goods
	if err := utils.Find(gc.db, &_goods, utils.NewWhereCond("owner_user", ownerUser), utils.WhereUIDCond(uuid)); err != nil {
		return model.GoodsFeild{}, err
	}

	field := model.GoodsFeild{
		GoodsName:   _goods.Name,
		GoodsTyp:    _goods.Typ,
		GoodsWeight: _goods.Weight,
	}

	// Store in hash
	data, _ := json.Marshal(field)
	gc.rdb.HSet(ctx, hashKey, uuid, data)
	gc.rdb.Expire(ctx, hashKey, 30*time.Minute)

	return field, nil
}

// Batchmodel.GoodsFeildSet with Hash
func (gc *goodsCache) BatchGoodsFeildSet(uuidList []string, ownerUser string) (map[string]model.GoodsFeild, error) {
	ctx := context.Background()
	hashKey := gc.goodsHashKey(ownerUser)
	result := make(map[string]model.GoodsFeild, len(uuidList))

	// Batch get from hash - HMGET for multiple fields
	cached, err := gc.rdb.HMGet(ctx, hashKey, uuidList...).Result()
	missingUUIDs := make([]string, 0)

	if err == nil {
		for i, val := range cached {
			if val != nil {
				var field model.GoodsFeild
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
		var _goodsList []model.Goods
		if err := utils.Find(gc.db, &_goodsList, utils.NewWhereCond("owner_user", ownerUser),
			utils.NewInCondFromString("uid", missingUUIDs)); err != nil {
			return result, err
		}

		// HMSET for batch storing in hash
		hashData := make(map[string]interface{})
		for _, goods := range _goodsList {
			field := model.GoodsFeild{
				GoodsName:   goods.Name,
				GoodsTyp:    goods.Typ,
				GoodsWeight: goods.Weight,
			}
			result[goods.UID] = field
			data, _ := json.Marshal(field)
			hashData[goods.UID] = data
		}

		if len(hashData) > 0 {
			gc.rdb.HMSet(ctx, hashKey, hashData)
			gc.rdb.Expire(ctx, hashKey, 30*time.Minute)
		}
	}

	return result, nil
}

// Cache invalidation with Hash
func (gc *goodsCache) InvalidateGoodsCache(uuid, ownerUser string) error {
	ctx := context.Background()
	hashKey := gc.goodsHashKey(ownerUser)
	return gc.rdb.HDel(ctx, hashKey, uuid).Err()
}

func (gc *goodsCache) InvalidateBatchGoodsCache(uuidList []string, ownerUser string) error {
	ctx := context.Background()
	hashKey := gc.goodsHashKey(ownerUser)
	return gc.rdb.HDel(ctx, hashKey, uuidList...).Err()
}

// Clear all cache for a specific user (useful when user deletes account)
func (gc *goodsCache) ClearUserGoodsCache(ownerUser string) error {
	ctx := context.Background()
	hashKey := gc.goodsHashKey(ownerUser)
	return gc.rdb.Del(ctx, hashKey).Err()
}

// CRUD operations with cache management
func (gc *goodsCache) CreateGoods(goods *model.Goods) error {
	if err := gc.db.Create(goods).Error; err != nil {
		return err
	}
	return gc.InvalidateGoodsCache(goods.UID, goods.OwnerUser)
}

func (gc *goodsCache) UpdateGoods(uuid, ownerUser string, updates map[string]interface{}) error {
	if err := gc.db.Model(&model.Goods{}).Where("uid = ? AND owner_user = ?", uuid, ownerUser).Updates(updates).Error; err != nil {
		return err
	}
	return gc.InvalidateGoodsCache(uuid, ownerUser)
}

func (gc *goodsCache) DeleteGoods(uuid, ownerUser string) error {
	if err := gc.db.Where("uid = ? AND owner_user = ?", uuid, ownerUser).Delete(&model.Goods{}).Error; err != nil {
		return err
	}
	return gc.InvalidateGoodsCache(uuid, ownerUser)
}
