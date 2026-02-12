package cron

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"go-additional-tools/ecron/mid"
	"go-additional-tools/ecron/model"
	"go-additional-tools/ecron/server/listen/service"
)

func getCronConfigByMsg(msg *nats.Msg) (model.CronJob, error) {
	var cronConf = model.CronJob{}
	err := gjson.Unmarshal(msg.Data, &cronConf)
	if err != nil {
		return cronConf, err
	}
	if cronConf.FuncName == "" {
		return cronConf, errors.New("数据解析异常")
	}

	return cronConf, nil
}

func AddCronSubscribe() error {
	return mid.Nats(model.CronNats).SubscribeRequest(model.CronAddNatsSubj, "", func(subj, queue string, msg *nats.Msg) (context.Context, []byte, error) {
		ctx := mid.SetMsgCtx(msg)
		g.Log().Info(ctx, subj, queue, "执行操作："+model.CronAddNatsSubj)
		var cronConf, err = getCronConfigByMsg(msg)
		if err != nil {
			return ctx, []byte(""), err
		}
		cron, err := service.GetCronService().AddCron(ctx, cronConf)
		return ctx, []byte(cron), err
	})
}
func DeleteCronSubscribe() error {
	return mid.Nats(model.CronNats).SubscribeRequest(model.CronDelNatsSubj, "", func(subj, queue string, msg *nats.Msg) (context.Context, []byte, error) {
		ctx := mid.SetMsgCtx(msg)
		g.Log().Info(ctx, subj, queue, "执行操作："+model.CronDelNatsSubj)
		var cronConf, err = getCronConfigByMsg(msg)
		if err != nil {
			return ctx, []byte(""), err
		}
		cron, err := service.GetCronService().DeleteCron(ctx, cronConf)
		return ctx, []byte(cron), err
	})
}

func StartCronSubscribe() error {
	return mid.Nats(model.CronNats).SubscribeRequest(model.CronStartNatsSubj, "", func(subj, queue string, msg *nats.Msg) (context.Context, []byte, error) {
		ctx := mid.SetMsgCtx(msg)
		g.Log().Info(ctx, subj, queue, "执行操作："+model.CronStartNatsSubj)
		var cronConf, err = getCronConfigByMsg(msg)
		if err != nil {
			return ctx, []byte(""), err
		}
		cron, err := service.GetCronService().StartCron(ctx, cronConf)
		return ctx, []byte(cron), err
	})
}
func StopCronSubscribe() error {
	return mid.Nats(model.CronNats).SubscribeRequest(model.CronStopNatsSubj, "", func(subj, queue string, msg *nats.Msg) (context.Context, []byte, error) {
		ctx := mid.SetMsgCtx(msg)
		g.Log().Info(ctx, subj, queue, "执行操作："+model.CronStopNatsSubj)
		var cronConf, err = getCronConfigByMsg(msg)
		if err != nil {
			return ctx, []byte(""), err
		}
		cron, err := service.GetCronService().StopCron(ctx, cronConf)
		return ctx, []byte(cron), err
	})
}

func InitCronSubscribe() error {
	err := AddCronSubscribe()
	if err != nil {
		return err
	}

	err = DeleteCronSubscribe()
	if err != nil {
		return err
	}
	err = StartCronSubscribe()
	if err != nil {
		return err
	}
	err = StopCronSubscribe()
	if err != nil {
		return err
	}

	return err

}
