package cache

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	SmsCodeExpire    = 5 * time.Minute
	SmsCountExpire   = 24 * time.Hour
	SmsMaxDailyCount = 5
)

type smsCache struct {
	rdb redis.UniversalClient
}

func NewSmsCache(rdb redis.UniversalClient) *smsCache {
	return &smsCache{rdb: rdb}
}

func (sc *smsCache) smsCodeKey(phone string) string {
	return fmt.Sprintf("%s:sms:code:%s", CachePrefix, phone)
}

func (sc *smsCache) smsCountKey(phone string) string {
	return fmt.Sprintf("%s:sms:count:%s", CachePrefix, phone)
}

func (sc *smsCache) SetVerifyCode(phone string) (string, error) {
	ctx := context.Background()

	// 检查今日发送次数
	over, err := sc.CheckTodayCount(phone)
	if err != nil {
		return "", err
	}
	if over {
		return "", fmt.Errorf("今日发送次数已达上限")
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 存储验证码
	codeKey := sc.smsCodeKey(phone)
	if err := sc.rdb.Set(ctx, codeKey, code, SmsCodeExpire).Err(); err != nil {
		return "", err
	}

	// 增加计数
	countKey := sc.smsCountKey(phone)
	pipe := sc.rdb.Pipeline()
	pipe.Incr(ctx, countKey)
	pipe.Expire(ctx, countKey, SmsCountExpire)
	_, err = pipe.Exec(ctx)

	return code, err
}

func (sc *smsCache) VerifyCode(phone, input string) (bool, error) {
	ctx := context.Background()
	codeKey := sc.smsCodeKey(phone)

	storedCode, err := sc.rdb.Get(ctx, codeKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if storedCode == input {
		sc.rdb.Del(ctx, codeKey) // 验证成功删除验证码
		return true, nil
	}

	return false, nil
}

func (sc *smsCache) CheckTodayCount(phone string) (bool, error) {
	ctx := context.Background()
	countKey := sc.smsCountKey(phone)

	count, err := sc.rdb.Get(ctx, countKey).Int64()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return count >= SmsMaxDailyCount, nil
}
