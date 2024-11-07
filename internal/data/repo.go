package data

import (
	"service-checker/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type repo struct {
	data *Data
	log  *log.Helper
}

// NewRepo .
func NewRepo(data *Data, logger log.Logger) biz.Repo {
	return &repo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
