package model

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type CronExec struct {
	ID          string `json:"id"`        // 主键ID
	CronExpr    string `json:"cron_expr"` // Cron表达式
	Params      string `json:"params"`    // 执行参数
	Type        CronSubj
	FuncName    string `json:"func_name"`   // 执行函数名
	Description string `json:"description"` // 任务描述
	Subject     string `json:"subject"`     // NATS主题/分组
	Timeout     string `json:"timeout"`     // 超时时间
	CheckLock   bool   `json:"check_lock"`  // 锁检查开关
	ProjName    string `json:"proj_name"`   // 工程名称
	ReTry       int    `json:"re_try"`
}

type CronExecOpt func(o *CronExec)

func WithId(id string) CronExecOpt {
	return func(o *CronExec) {
		o.ID = id
	}
}

func WithCronExpr(cronExpr string) CronExecOpt {
	return func(o *CronExec) {
		o.CronExpr = cronExpr
	}
}

func WithParams(params string) CronExecOpt {
	return func(o *CronExec) {
		o.Params = params
	}
}
func WithReTry(reRry int) CronExecOpt {
	return func(o *CronExec) {
		o.ReTry = reRry
	}
}
func WithFuncName(funcName string) CronExecOpt {
	return func(o *CronExec) {
		o.FuncName = funcName
	}
}
func WithDescription(description string) CronExecOpt {
	return func(o *CronExec) {
		o.Description = description
	}
}

func WithSubject(subject string) CronExecOpt {
	return func(o *CronExec) {
		o.Subject = subject
	}
}

func WithTimeout(timeout string) CronExecOpt {
	return func(o *CronExec) {
		o.Timeout = timeout
	}
}

func WithCheckLock(checkLock bool) CronExecOpt {
	return func(o *CronExec) {
		o.CheckLock = checkLock
	}
}

func WithType(t CronSubj) CronExecOpt {
	return func(o *CronExec) {
		o.Type = t
	}
}
func NewCronExec(ops ...CronExecOpt) CronExec {

	var c = CronExec{
		ID:   gconv.String(gtime.TimestampNano()),
		Type: CronAddNatsSubjType,
	}

	for _, op := range ops {
		op(&c)
	}

	return c
}
