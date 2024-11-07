package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/nuominmin/notifier"
	workerreloader "github.com/nuominmin/worker-reloader"
	"service-checker/internal/conf"
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
		interval := s.c.Services[i].Interval
		if interval <= 0 {
			interval = 60
		}
		s.workerPool.Start(s.c.Services[i].Name, s.healthy(&Checker{
			Name:       s.c.Services[i].Name,
			URL:        s.c.Services[i].Url,
			MaxRetries: s.c.Services[i].MaxRetries,
		}), time.Duration(interval)*time.Second)
	}
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	log.Info("stopping service")
	return nil
}
