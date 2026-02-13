package client

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/nats-io/nats.go"
	"go-additional-tools/ecron/mid"
	"go-additional-tools/ecron/model"
	"sync"
	"time"
)

var (
	lockMap  = sync.Map{}
	startMap = sync.Map{}
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

// 定时任务执行的订阅
func ExecSubscribe(subj string, name string) error {
	if subj == "" {
		return errors.New("订阅主题不能为空 cron.subjectName")
	}

	queue := "ExecCronRequest" // 固定参数，方便单一订阅
	_, alreadyExists := startMap.LoadOrStore(subj, struct{}{})
	if alreadyExists {
		return nil
	}

	// 执行定时任务
	err := mid.Nats(model.ExecNats).SubscribeRequest(subj, queue, func(subj, queue string, msg *nats.Msg) (context.Context, []byte, error) {
		msg.Header.Set(mid.RespondConsumerKey, name)
		ctx := mid.SetMsgCtx(msg)
		var cronConf, err = getCronConfigByMsg(msg)
		lockMap.Store(cronConf.Subject+":"+cronConf.FuncName+":"+cronConf.ID, struct{}{})
		if err != nil {
			return ctx, []byte(""), err
		}
		execService, err := GetCronExecService(cronConf.FuncName)
		if err != nil {
			lockMap.Delete(cronConf.Subject + ":" + cronConf.FuncName + ":" + cronConf.ID)
			return ctx, []byte(""), err
		}
		var (
			fh              string
			timectx, cannel = context.WithCancel(context.Background())
		)
		go func() {
			defer func() {
				lockMap.Delete(cronConf.Subject + ":" + cronConf.FuncName + ":" + cronConf.ID)
				cannel()
			}()
			fh, err = execService(ctx, cronConf)
		}()
		select {
		case <-time.After(time.Second):
			return ctx, []byte(fh), errors.New("执行中>>>")
		case <-timectx.Done():
			return ctx, []byte(fh), err
		}

	})

	if err != nil {
		return err
	}
	// 用于校验任务是否运行
	err = mid.Nats(model.ExecNats).Subscribe(model.CronCheckLockNatsSubj+"."+subj, guid.S(), func(subj, queue string, msg *nats.Msg) {
		msg.Header.Set(mid.RespondConsumerKey, name)
		ctx := mid.SetMsgCtx(msg)
		var cronConf, err = getCronConfigByMsg(msg)
		if err != nil {
			g.Log().Error(ctx, "cronCheckLockNatsSubj:", err)
			return
		}

		_, ok := lockMap.Load(cronConf.Subject + ":" + cronConf.FuncName + ":" + cronConf.ID)
		if ok {
			//g.Log().Info(ctx, "cronCheckRespond")
			msg.Data = []byte(model.CronCheckIsLock)
			err = msg.RespondMsg(msg)
		} else {
			msg.Data = []byte(model.CronCheckNoLock)
			err = msg.RespondMsg(msg)
		}
		if err != nil {
			g.Log().Error(ctx, "cronCheckRespond:", err)
			return
		}
		//g.Log().Info(ctx, "checkLockFinish:"+cronConf.CSubj+":"+cronConf.CFunc+":"+cronConf.CKey)
	})
	if err != nil {
		mid.Nats(model.ExecNats).CloseSubscribe(subj, queue)
	}

	return err

}

func AddCronRequest(ctx context.Context, p model.CronExec, f ExecFunc) (string, error) {
	if p.ID == "" {
		p.ID = guid.S()
	}
	if p.CronExpr == "" {
		return "", errors.New("定时任务不能为空")
	}
	//if p.Subject == "" {
	//	return "", errors.New("订阅消息名称不能为空")
	//}

	if p.FuncName == "" {
		return "", errors.New("执行函数名称不能为空")
	}
	get, err2 := g.Cfg().Get(gctx.GetInitCtx(), "cron.subjectName")
	if err2 == nil && get.String() != "" {
		p.Subject = get.String()
	}

	if p.Subject == "" {
		return "", errors.New("未设置定时任务订阅主题 cron.subjectName")
	}

	get, err2 = g.Cfg().Get(gctx.GetInitCtx(), "cron.projectName")
	if err2 == nil && get.String() != "" {
		p.ProjName = get.String()
	}
	if p.ProjName == "" {
		return "", errors.New("未设置定时来源名称 cron.projectName")
	}
	ips, err2 := gipv4.GetIntranetIp()
	if err2 == nil && ips != "" {
		p.ProjName = p.ProjName + "(" + ips + ")"
	}

	RegisterIExecService(p.FuncName, f)
	err := ExecSubscribe(p.Subject, p.ProjName)
	if err != nil {
		return "", err
	}

	modl := model.CronJob{
		ID:          p.ID,
		FuncName:    p.FuncName,
		Description: p.Description,
		CronExpr:    p.CronExpr,
		Subject:     p.Subject,
		LastRun:     "",
		Status:      0,
		ErrorCode:   0,
		Msg:         "",
		Params:      p.Params,
		Timeout:     p.Timeout,
		CheckLock:   p.CheckLock,
		ReTry:       p.ReTry,
		ProjName:    p.ProjName,
	}
	var sujson = gjson.New(modl).MustToJson()

	request, err := mid.Nats(model.ExecNats).Request(ctx, string(p.Type), sujson, 10*time.Second)
	if err != nil {
		return "", err
	}
	msg := mid.GetErrorMsg(request)
	if msg != "" {
		return "", errors.New(msg)
	}
	return string(request.Data), err

}
