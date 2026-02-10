package utils

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var renewScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0
`)

type LeaderElector struct {
	rdb          *redis.Client
	instanceID   string
	leaderKey    string
	leaderTTL    time.Duration
	mu           sync.Mutex
	isLeader     bool
	leaderCancel context.CancelFunc
}

func NewLeaderElector(rdb *redis.Client, instanceID string, leaderKey string, leaderTTL time.Duration) *LeaderElector {
	return &LeaderElector{
		rdb:        rdb,
		instanceID: instanceID,
		leaderKey:  leaderKey,
		leaderTTL:  leaderTTL,
	}
}

func (l *LeaderElector) tryAcquire(ctx context.Context) (bool, error) {
	return l.rdb.SetNX(ctx, l.leaderKey, l.instanceID, l.leaderTTL).Result()
}

func (l *LeaderElector) renew(ctx context.Context) (bool, error) {
	res, err := renewScript.Run(ctx, l.rdb, []string{l.leaderKey}, l.instanceID, int(l.leaderTTL.Milliseconds())).Int()
	return res == 1, err
}

func (l *LeaderElector) startLeader(ctx context.Context, onLeader func(ctx context.Context)) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.isLeader {
		return
	}

	leaderCtx, cancel := context.WithCancel(ctx)
	l.isLeader = true
	l.leaderCancel = cancel

	go onLeader(leaderCtx)
}

func (l *LeaderElector) stopLeader() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isLeader {
		return
	}

	l.isLeader = false
	if l.leaderCancel != nil {
		l.leaderCancel()
		l.leaderCancel = nil
	}
}

func (l *LeaderElector) Run(ctx context.Context, onLeader func(ctx context.Context)) {
	go func() {
		ticker := time.NewTicker(l.leaderTTL / 2)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if l.isLeader {
					ok, err := l.renew(ctx)
					if err != nil || !ok {
						l.stopLeader()
					}
				} else {
					ok, err := l.tryAcquire(ctx)
					if ok && err == nil {
						l.startLeader(ctx, onLeader)
					}
				}
			case <-ctx.Done():
				l.stopLeader()
				return
			}
		}
	}()
}
