package service

import (
	"context"
	"go-additional-tools/ecron/model"
)

var localICronService ICronService = nil

type ICronService interface {
	AddCron(ctx context.Context, m model.CronJob) (string, error)
	InsertCronLog(ctx context.Context, m model.CronJob) (string, error)
	DeleteCron(ctx context.Context, m model.CronJob) (string, error)
	StartCron(ctx context.Context, m model.CronJob) (string, error)
	StopCron(ctx context.Context, m model.CronJob) (string, error)
	ListCron(ctx context.Context) ([]model.CronJob, error)
	GetCronOne(ctx context.Context, id string) (*model.CronJob, error)
}

func RegisterCronService(i ICronService) {
	localICronService = i
}

func GetCronService() ICronService {
	if localICronService == nil {
		panic("CronService not register")
	}
	return localICronService
}
