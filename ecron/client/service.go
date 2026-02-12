package client

import (
	"context"
	"errors"
	"go-additional-tools/ecron/model"
	"sync"
)

type ExecFunc func(ctx context.Context, model model.CronJob) (string, error)

var (
	localIExecServiceMap = make(map[string]ExecFunc)
	lock                 sync.Mutex
)

func RegisterIExecService(name string, f ExecFunc) {
	lock.Lock()
	defer lock.Unlock()
	localIExecServiceMap[name] = f
}

func GetCronExecService(name string) (ExecFunc, error) {
	lock.Lock()
	defer lock.Unlock()

	if i, ok := localIExecServiceMap[name]; ok {
		return i, nil
	} else {
		return nil, errors.New("not found service:" + name)
	}

}
