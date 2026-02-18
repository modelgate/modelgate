package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
	"github.com/modelgate/modelgate/pkg/utils"
)

func (s *Service) StartWorker(ctx context.Context) {
	log.Info("Start worker...")
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info("leader lost, stop worker")
				return
			case <-ticker.C:
				if err := s.saveCachedUsageToDB(ctx); err != nil {
					log.Error("save cached usage to db error", err)
				}
				if err := s.cleanExpiredUsageKeys(ctx); err != nil && !errors.Is(err, redis.Nil) {
					log.Error("clean expired usage keys error", err)
				}
			}
		}
	}()
}

func (s *Service) saveCachedUsageToDB(ctx context.Context) (err error) {
	var keys []string
	for _, prefix := range model.AllUsagePrefixs {
		zkey := s.getZSetKey(prefix)
		keys, err = s.redisClient.ZRangeByScore(ctx, zkey, &redis.ZRangeBy{
			Min: "-inf",
			Max: "+inf",
		}).Result()
		if err != nil {
			return
		}
		for _, key := range keys {
			if strings.HasPrefix(key, model.UsageProviderPrefix) {
				err = s.saveRequestUsageToDB(ctx, key)
			} else if strings.HasPrefix(key, model.UsageAccountApiKeyPrefix) || strings.HasPrefix(key, model.UsageProviderApiKeyPrefix) {
				err = s.saveApiKeyUsageToDB(ctx, key)
			}
		}
	}
	return
}

var delScript = redis.NewScript(`
	redis.call('DEL', KEYS[1])
	redis.call('ZREM', KEYS[2], KEYS[1])
	return
`)

func (s *Service) cleanExpiredUsageKeys(ctx context.Context) (err error) {
	var keys []string
	for _, prefix := range model.AllUsagePrefixs {
		zkey := s.getZSetKey(prefix)
		keys, err = s.redisClient.ZRangeByScore(ctx, zkey, &redis.ZRangeBy{
			Min: "-inf",
			Max: strconv.FormatInt(time.Now().Unix(), 10),
		}).Result()
		if err != nil {
			return
		}
		for _, key := range keys {
			err = delScript.Run(ctx, s.redisClient, []string{key, zkey}).Err()
			if err != nil {
				return
			}
		}
	}
	return
}

// key format: usage:request:{provider_code}:{metric}
func (s *Service) saveRequestUsageToDB(ctx context.Context, key string) (err error) {
	parts := strings.Split(key, ":")
	if len(parts) < 4 {
		err = errors.New("invalid key")
		return
	}

	providerCode := parts[2]
	metric := parts[3]
	statTime, _ := utils.ParseDate(utils.FormatTime(time.Now(), "YmdH"), "YmdH", time.Local)

	err = s.processingCache(ctx, key, func(ctx context.Context, value int64) (err error) {
		usage := &model.RelayHourlyUsage{
			Time:         statTime,
			ProviderCode: providerCode,
		}
		switch metric {
		case model.MetricTotal:
			usage.TotalRequest = value
			_, err = s.relayUsageDao.Update(ctx, &model.RelayUsageFilter{ID: db.Eq(int64(1))}, map[string]any{
				"total_request": gorm.Expr("total_request + ?", value),
			})
			if err != nil {
				return
			}
		case model.MetricSuccess:
			usage.TotalSuccess = value
		case model.MetricFailed:
			usage.TotalFailed = value
		case model.MetricUsage:
			usage.TotalPoint = value
			_, err = s.relayUsageDao.Update(ctx, &model.RelayUsageFilter{ID: db.Eq(int64(1))}, map[string]any{
				"total_point": gorm.Expr("total_point + ?", value),
			})
			if err != nil {
				return
			}
		}
		err = s.relayHourlyUsageDao.Create(ctx, usage)
		return
	})
	return
}

// key format: usage:account_api_key:{api_key_id}:usage or usage:provider_api_key:{api_key_id}:usage
func (s *Service) saveApiKeyUsageToDB(ctx context.Context, key string) (err error) {
	parts := strings.Split(key, ":")
	if len(parts) < 4 {
		err = errors.New("invalid key")
		return
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		err = errors.New("invalid key")
		return
	}
	err = s.processingCache(ctx, key, func(ctx context.Context, value int64) (err error) {
		switch {
		case strings.HasPrefix(key, model.UsageAccountApiKeyPrefix):
			_, err = s.accountApiKeyDao.Update(ctx, &model.AccountApiKeyFilter{ID: db.Eq(id)},
				map[string]any{
					"last_used_at": time.Now(),
					"quote_used":   gorm.Expr("quote_used + ?", value),
				})
		case strings.HasPrefix(key, model.UsageProviderApiKeyPrefix):
			_, err = s.providerApiKeyDao.Update(ctx, &model.ProviderApiKeyFilter{ID: db.Eq(id)},
				map[string]any{
					"last_used_at": time.Now(),
					"quote_used":   gorm.Expr("quote_used + ?", value),
				})
		}
		return
	})
	return
}

func (s *Service) processingCache(ctx context.Context, key string, f func(ctx context.Context, value int64) error) (err error) {
	delta, err := s.redisClient.GetSet(ctx, key, 0).Int64()
	if err == redis.Nil {
		err = nil
	} else if err != nil {
		return
	}
	if delta == 0 {
		return
	}
	err = f(ctx, delta)
	return
}

var incrWithZAddScript = redis.NewScript(`
	local key = KEYS[1]
	local zset_key = KEYS[2]
	local value = ARGV[1]
	local current_time = tonumber(ARGV[2])
	local expire_offset = tonumber(ARGV[3]) or 3600
	local expire_time = current_time + expire_offset

	redis.call('INCRBY', key, value)
	redis.call('ZADD', zset_key, expire_time, key)
	return expire_time
`)

func (s *Service) incrMetricValue(ctx context.Context, prefix string, objectMetric string, value int64, expireOffset int64) (expireTime int64, err error) {
	key := prefix + objectMetric
	zkey := s.getZSetKey(prefix)
	expireTime, err = incrWithZAddScript.Run(ctx, s.redisClient, []string{key, zkey}, value, time.Now().Unix(), expireOffset).Int64()
	return
}

func (s *Service) getZSetKey(prefix string) string {
	return prefix + "dirty_keys"
}

// AddRequestUsage adds request usage to redis
func (s *Service) AddRequestUsage(ctx context.Context, providerCode string, metric string, value int64) (err error) {
	_, err = s.incrMetricValue(ctx, model.UsageProviderPrefix, fmt.Sprintf("%s:%s", providerCode, metric), value, 3600)
	return
}

// AddPointUsage adds point usage to redis
func (s *Service) AddPointUsage(ctx context.Context, providerCode string, providerApiKeyId, accountApiKeyId int64, value int64) (err error) {
	_, err = s.incrMetricValue(ctx, model.UsageProviderPrefix, fmt.Sprintf("%s:%s", providerCode, model.MetricUsage), value, 3600)
	if err != nil {
		return
	}
	_, err = s.incrMetricValue(ctx, model.UsageProviderApiKeyPrefix, fmt.Sprintf("%d:%s", providerApiKeyId, model.MetricUsage), value, 3600)
	if err != nil {
		return
	}
	_, err = s.incrMetricValue(ctx, model.UsageAccountApiKeyPrefix, fmt.Sprintf("%d:%s", accountApiKeyId, model.MetricUsage), value, 3600)
	return
}
