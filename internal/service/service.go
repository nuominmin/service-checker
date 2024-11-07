package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
	pb "service-checker/api"
	"service-checker/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService)

// Service is a service.
type Service struct {
	pb.UnimplementedV1Server

	svc *biz.Service
}

// NewService new a service.
func NewService(svc *biz.Service) (s *Service, cf func(), err error) {
	s = &Service{
		svc: svc,
	}

	return s, s.Close, nil
}

// Close close the resource.
func (s *Service) Close() {
	log.Info("close service resource. ")
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// List .
func (s *Service) List(ctx context.Context, in *pb.ListReq) (*pb.ListResp, error) {
	return &pb.ListResp{
		Total:    0,
		DataList: make([]*pb.ListResp_Data, 0),
	}, nil
}
