package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CheckReady 检查服务是否 ready
func CheckReady(i do.Injector) {
	if err := checkRedisReady(i); err != nil {
		log.Fatalf("redis ping err: %v", err)
	}
	if err := checkDatabaseReady(i); err != nil {
		log.Fatalf("database ping err: %v", err)
	}
}

func checkRedisReady(i do.Injector) error {
	redisClient := do.MustInvoke[*redis.Client](i)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return redisClient.Ping(ctx).Err()
}

func checkDatabaseReady(i do.Injector) error {
	dbConn := do.MustInvoke[*gorm.DB](i)
	dbc, err := dbConn.DB()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := dbc.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
