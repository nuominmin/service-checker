package factory

import (
	"github.com/google/wire"
	"github.com/nuominmin/cache"
	"github.com/nuominmin/notifier"
	"github.com/nuominmin/notifier/lark"
	workerreloader "github.com/nuominmin/worker-reloader"
	"service-checker/internal/conf"
)

var ProviderSet = wire.NewSet(
	NewAlert,
	NewCache,
	NewNewWorkerPool,
)

func NewAlert(data *conf.Data) (notifier.Notifier, error) {
	alert := lark.NewNotifier(data.GetAlertTokens()...)
	if env := data.GetEnv(); env != "" {
		alert.SetIdentity(env)
	}
	return alert, nil
}

func NewCache() cache.Cache {
	return cache.NewCache()
}

func NewNewWorkerPool() (workerreloader.WorkerPoolManager, func()) {
	wpm := workerreloader.NewWorkerPool()
	return wpm, wpm.StopAll
}
