package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/nuominmin/notifier"
	workerreloader "github.com/nuominmin/worker-reloader"
	"service-checker/internal/conf"
	"sync"
	"time"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewService)

// Service .
type Service struct {
	alert      notifier.Notifier
	workerPool workerreloader.WorkerPoolManager
	c          *conf.Data
	log        *log.Helper

	checkers   []*Checker   // 存储所有的 Checker
	checkersMu sync.RWMutex // 用于并发访问 checkers 的读写锁
}

// NewService .
func NewService(c *conf.Data, alert notifier.Notifier, workerPool workerreloader.WorkerPoolManager, logger log.Logger) *Service {
	return &Service{
		alert:      alert,
		workerPool: workerPool,
		c:          c,
		log:        log.NewHelper(logger),
	}
}

func (s *Service) Start(ctx context.Context) error {
	log.Info("starting service")
	for i := 0; i < len(s.c.Services); i++ {
		// 创建 Checker 实例
		checker := &Checker{
			Name: s.c.Services[i].Name,
			URL:  s.c.Services[i].Url,
		}

		// 将 Checker 添加到列表中
		s.addChecker(checker)
		s.workerPool.Start(s.c.Services[i].Name, s.healthy(checker), time.Duration(1)*time.Second)
	}
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	log.Info("stopping service")
	return nil
}
