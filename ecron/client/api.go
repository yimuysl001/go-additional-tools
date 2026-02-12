package client

import (
	"context"
	"go-additional-tools/ecron/model"
)

func AddCron(ctx context.Context, p model.CronExec, f ExecFunc) (string, error) {
	p.Type = model.CronStartNatsSubj
	return AddCronRequest(ctx, p, f)
}

func AddSingleton(ctx context.Context, p model.CronExec, f ExecFunc) (string, error) {
	p.Type = model.CronStartNatsSubj
	p.CheckLock = true
	return AddCronRequest(ctx, p, f)
}
