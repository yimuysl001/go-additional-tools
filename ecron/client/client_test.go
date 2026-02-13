package client

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"go-additional-tools/econf"
	"go-additional-tools/ecron/model"
	"testing"
	"time"
)

func TestCli(t *testing.T) {
	econf.MustInitConf()
	_, err := AddSingleton(gctx.GetInitCtx(), model.CronExec{
		ID:          "TestCronFuncaaa",
		CronExpr:    "0/30 * * * * *",
		Params:      "测试",
		FuncName:    "CronFunc",
		Description: "数据测试",
		Timeout:     "10s",
		CheckLock:   true,
		ReTry:       3,
	}, func(ctx context.Context, p model.CronJob) (string, error) {
		g.Log().Info(ctx, "测试执行开始")
		//time.Sleep(time.Minute)
		g.Log().Info(ctx, "测试执行结束")
		return "ok", nil
	})
	if err != nil {
		g.Log().Error(gctx.New(), err)
	} else {
		g.Log().Info(gctx.New(), "添加成功")
	}
	time.Sleep(time.Hour)
}
