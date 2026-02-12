package impl

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
	"go-additional-tools/ecron/model"
	"go-additional-tools/ecron/server/listen/cron"
	sysmnt "go-additional-tools/ecron/server/listen/service/impl/glc"
)

type SysNmtDb struct {
	Name string
}

func NewSysNmtDb(name string) *SysNmtDb {
	if name == "" {
		name = "sysmnt"
	}
	return &SysNmtDb{
		Name: name,
	}
}

func (c *SysNmtDb) AddCron(ctx context.Context, m model.CronJob) (string, error) {
	storage := sysmnt.NewSysmntStorage(ctx, c.Name)
	return m.GetKey(), storage.SaveModel(ctx, &m)
}

func (c *SysNmtDb) InsertCronLog(ctx context.Context, m model.CronJob) (string, error) {

	return c.AddCron(ctx, m)
}

func (c *SysNmtDb) DeleteCron(ctx context.Context, m model.CronJob) (string, error) {

	storage := sysmnt.NewSysmntStorage(ctx, c.Name)
	err := storage.DelModel(ctx, &m)
	if err != nil {
		return "", err
	}
	gcron.Stop(m.GetKey())
	gcron.Remove(m.GetKey())
	return m.GetKey(), err
}

func (c *SysNmtDb) StartCron(ctx context.Context, m model.CronJob) (string, error) {
	m.Status = 1
	cronkey, err := c.AddCron(ctx, m)
	if err != nil {
		return "", err
	}
	gcron.Stop(cronkey)
	gcron.Remove(cronkey)
	return cronkey, cron.StartCron(ctx, m)

}

func (c *SysNmtDb) StopCron(ctx context.Context, m model.CronJob) (string, error) {
	m.Status = 2
	cronkey, err := c.AddCron(ctx, m)
	if err != nil {
		return "", err
	}
	gcron.Stop(cronkey)
	gcron.Remove(cronkey)
	return cronkey, nil
}

func (c *SysNmtDb) GetCronOne(ctx context.Context, id string) (*model.CronJob, error) {
	storage := sysmnt.NewSysmntStorage(ctx, c.Name)

	m := &model.CronJob{
		ID: "~KEY:" + id,
	}

	err, _ := storage.GetModel(ctx, m)
	if err != nil {
		return nil, err
	}

	return m, err
}

func (c *SysNmtDb) ListCron(ctx context.Context) ([]model.CronJob, error) {

	storage := sysmnt.NewSysmntStorage(ctx, c.Name)

	job := &model.CronJob{}
	keys := storage.GetAllKeys(ctx, job.GetPrefix())

	var list = make([]model.CronJob, len(keys))
	for i, key := range keys {
		one, err := c.GetCronOne(ctx, key)
		if err != nil {
			return nil, err
		}
		list[i] = *one
	}

	return list, nil
}
