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

func (s *Service) Services(ctx context.Context, in *empty.Empty) (*pb.ServicesResp, error) {
	allCheckers := s.svc.GetAllCheckers()
	services := make([]*pb.ServicesResp_Service, len(allCheckers))
	for i := 0; i < len(allCheckers); i++ {
		services[i] = &pb.ServicesResp_Service{
			Name:   allCheckers[i].Name,
			Status: uint32(allCheckers[i].Status),
		}
	}
	return &pb.ServicesResp{
		Services: services,
	}, nil
}
