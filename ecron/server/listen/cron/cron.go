package cron

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nats-io/nats.go"
	"go-additional-tools/ecron/mid"
	"go-additional-tools/ecron/model"
	"go-additional-tools/ecron/server/listen/service"
	"strings"
	"time"
)

func getSubjCount(ctx context.Context, subj string) int {
	var query = []byte(`{"subscriptions":true,"filter_subject":"` + subj + `"}`)
	request, err := mid.Nats(model.ExecNats).Request(ctx, "$SYS.REQ.SERVER.PING.CONNZ", query, 2*time.Second)
	if err != nil {
		g.Log().Debug(ctx, "mq.Nats:"+err.Error())
		return 0
	}
	return gjson.New(request.Data).Get("data.num_connections").Int()

}

func ExecRequest(ctx context.Context, m *model.CronJob) (string, error) {

	json, err := gjson.New(m).ToJson()
	if err != nil {
		return "", err
	}
	timed := 5 * time.Second
	if m.Timeout != "" {
		timed, err = gtime.ParseDuration(m.Timeout)
	}
	if err != nil {
		return "", err
	}

	if m.CheckLock {
		g.Log().Info(ctx, "校验是否完成执行")
		count := getSubjCount(ctx, model.CronCheckLockNatsSubj+"."+m.Subject)
		if count == 0 {
			return "", nats.ErrNoResponders
		}
		g.Log().Info(ctx, "共", count, "个消费者")
		err = mid.Nats(model.ExecNats).PublishRequest(ctx, model.CronCheckLockNatsSubj+"."+m.Subject, json, count, func(ctx context.Context, msg *nats.Msg) error {
			name := mid.GetProjName(msg)
			if name != "" {
				m.ProjName = name
				g.Log().Info(ctx, "执行工程名称：", name)
			}
			g.Log().Debug(ctx, "get subscribe:"+string(msg.Data))
			if strings.EqualFold(model.CronCheckIsLock, string(msg.Data)) {
				return errors.New("存在订阅未完成")
			}
			return nil
		})
		if err != nil {
			return "", err
		}

	}
	request, err := mid.Nats(model.ExecNats).Request(ctx, m.Subject, json, timed)
	if err != nil {
		return "", err
	}
	msg := mid.GetErrorMsg(request)
	name := mid.GetProjName(request)
	if name != "" {
		m.ProjName = name
		g.Log().Info(ctx, "执行工程名称：", name)
	}

	if msg != "" {
		return "", errors.New(msg)
	}
	return string(request.Data), err

}

//func StartCron(ctx context.Context, config model.CronJob) error {
//	var ckey = config.Subject + ":" + config.FuncName + ":" + config.ID
//	g.Log().Info(ctx, "启动定时任务：", ckey)
//	_, err := gcron.AddSingleton(gctx.New(), config.CronExpr, func(ctx context.Context) {
//		config.ChainId = gctx.CtxId(ctx)
//		g.Log().Info(ctx, ckey+":"+config.Description+":开始执行")
//		config.Status = 1
//		if config.RunStatus != 4 {
//			config.LastRun = time.Now().Format("2006-01-02 15:04:15")
//			config.RunStatus = 1
//		}
//
//		defer func() {
//			g.Log().Info(ctx, ckey+":"+config.Description+":执行完成")
//			if config.RunStatus == 4 {
//				config.LastEnd = ""
//
//			} else if config.RunStatus == 3 {
//				config.Status = 2
//			} else {
//				config.RunStatus = 2
//			}
//
//			service.GetCronService().InsertCronLog(ctx, config)
//		}()
//		request, err2 := ExecRequest(ctx, &config)
//		config.LastEnd = time.Now().Format("2006-01-02 15:04:15")
//		config.Msg = request
//		if err2 != nil {
//			g.Log().Info(ctx, ckey+":"+config.Description+":", err2)
//			if errors.Is(err2, nats.ErrNoResponders) {
//				g.Log().Info(ctx, ckey+":"+config.Description+":无订阅者，停止定时任务")
//				config.RunStatus = 2
//				config.ErrorCode = 0
//				config.Msg = "无订阅者，停止定时任务"
//				service.GetCronService().StopCron(ctx, config)
//			} else if strings.Contains(err2.Error(), "存在订阅未完成") {
//				config.RunStatus = 4
//				config.Msg = "存在订阅未完成"
//			} else {
//				config.RunStatus = 3
//				config.ErrorCode = 1
//				config.Msg = err2.Error()
//			}
//			return
//		}
//		config.RunStatus = 2
//
//		g.Log().Info(ctx, ckey+":"+config.Description+":", request)
//
//	}, ckey)
//	return err
//
//}

func StartCron(ctxN context.Context, config model.CronJob) error {
	const TimeLayout = "2006-01-02 15:04:15"
	var ckey = config.Subject + ":" + config.FuncName + ":" + config.ID

	g.Log().Info(ctxN, "启动定时任务：", ckey)

	jobConfig := config // 局部副本避免并发冲突
	ctx := gctx.New()
	jobConfig.ChainId = gctx.CtxId(ctx)
	jobConfig.ReCount = 0
	_, err := gcron.AddSingleton(gctx.New(), jobConfig.CronExpr, func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Error(ctx, "panic occurred:", r)
				jobConfig.RunStatus = 3
				jobConfig.Msg = fmt.Sprintf("panic: %v", r)
			}

			if jobConfig.RunStatus != 4 {
				jobConfig.LastEnd = time.Now().Format(TimeLayout)
				//go func() {
				service.GetCronService().InsertCronLog(ctx, jobConfig)
				//}()
			}
		}()

		jobConfig.Status = 1
		if jobConfig.RunStatus != 4 {
			jobConfig.LastRun = time.Now().Format(TimeLayout)
			jobConfig.RunStatus = 1
			service.GetCronService().InsertCronLog(ctx, jobConfig)
		}

		request, err2 := ExecRequest(ctx, &jobConfig)

		jobConfig.Msg = request
		if err2 != nil {
			updateConfigOnError(ctx, &jobConfig, err2)
			return
		}
		config.ReCount = 0
		jobConfig.RunStatus = 2

	}, ckey)

	return err
}

func updateConfigOnError(ctx context.Context, config *model.CronJob, err error) {

	if errors.Is(err, nats.ErrNoResponders) {
		config.ErrorCode = 0
		config.Msg = "无订阅者，停止定时任务：重试第" + gconv.String(config.ReCount) + "次"
		if config.ReCount > config.ReTry {
			config.RunStatus = 2
			service.GetCronService().StopCron(ctx, *config)
		}
		config.ReCount = config.ReCount + 1
	} else if strings.Contains(err.Error(), "存在订阅未完成") {
		config.RunStatus = 4
		config.Msg = "存在订阅未完成"
		config.ReCount = 0
	} else {
		config.RunStatus = 3
		config.ErrorCode = 1
		config.Msg = err.Error()
		config.ReCount = 0
	}
}

func StartInitCron() error {
	err := InitCronSubscribe()
	if err != nil {
		return err
	}

	cron, err := service.GetCronService().ListCron(gctx.GetInitCtx())
	if err != nil {
		return err
	}
	for _, config := range cron {
		if config.Status != 1 {
			//if len(cron) > 50 {
			//	service.GetCronService().DeleteCron(gctx.GetInitCtx(), config)
			//}
			continue
		}

		err = StartCron(gctx.New(), config)
		if err != nil {
			g.Log().Error(gctx.GetInitCtx(), "gcron.AddSingleton:"+err.Error())
			break
		}

	}

	if err != nil {
		gcron.Stop()
	}

	return err

}
