package main

import (
	"context"
	"template/pkg/services/scan"
)

type Service struct {
	scan       *scan.Engine
	ctx        context.Context
	shutdownFn context.CancelFunc
}

func NewService(opt Option) *Service {
	rootCtx, shutdownFn := context.WithCancel(context.Background())

	scan, err := scan.NewEngine(rootCtx)
	if err != nil {
		panic(err)
	}

	service := &Service{
		scan:       scan,
		ctx:        rootCtx,
		shutdownFn: shutdownFn,
	}

	return service
}

func (service *Service) Init() {
	go service.scan.Run(service.ctx)
	//go service.scan.Run2(service.ctx)
}

func (service *Service) UnInit() {
	service.shutdownFn()
}
