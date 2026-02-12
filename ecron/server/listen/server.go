package listen

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"go-additional-tools/ecron/server/listen/cron"
	"go-additional-tools/ecron/server/listen/service"
	"go-additional-tools/ecron/server/listen/service/impl"
	"go-additional-tools/ecron/server/nats"
	"go-additional-tools/ecron/server/web"
	"time"
)

func StartSubjCron() error {
	get, err := g.Cfg().Get(gctx.GetInitCtx(), "cron.subscribe.dbpath")
	if err != nil {
		g.Log().Debug(gctx.GetInitCtx(), "get subscribe db error:"+err.Error())
		return err
	}
	if get.String() == "" {
		g.Log().Info(gctx.GetInitCtx(), "subscribe db is null")
		return err
	}
	service.RegisterCronService(impl.NewSysNmtDb("cron"))

	go func() {
		web.CronWeb()

	}()

	//db := get.String()
	//if db == "" {
	//	g.Log().Info(gctx.GetInitCtx(), "subscribe db is null")
	//	return errors.New("subscribe db is null")
	//}

	//cronsub, err := impl.NewLevelDBCron(db)
	//if err != nil {
	//	g.Log().Error(gctx.GetInitCtx(), "new leveldb cron error:"+err.Error())
	//	return err
	//}

	return cron.StartInitCron()

}

func StartNatsAndCron() {
	go func() {
		g.Log().Debug(gctx.GetInitCtx(), "启动定时")
		time.Sleep(time.Second)
		err2 := StartSubjCron()
		if err2 != nil {
			g.Log().Debug(gctx.GetInitCtx(), "定时启动失败:"+err2.Error())
		} else {
			g.Log().Debug(gctx.GetInitCtx(), "定时启动完成")
		}

	}()
	nats.NatsServer()
}
