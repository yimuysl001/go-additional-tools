package model

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"strings"
)

type CronJob struct {
	ID          string `json:"id"`          // 主键ID
	FuncName    string `json:"funcName"`    // 执行函数名
	Description string `json:"description"` // 任务描述
	CronExpr    string `json:"cronExpr"`    // Cron表达式
	Subject     string `json:"subject"`     // NATS主题/分组
	Params      string `json:"params"`      // 执行参数
	Timeout     string `json:"timeout"`     // 超时时间
	CheckLock   bool   `json:"checkLock"`   // 锁检查开关

	LastRun   string `json:"lastRun"`   // 最后执行时间
	LastEnd   string `json:"lastEnd"`   // 执行完成时间
	Status    int    `json:"status"`    // 任务状态 1 启用 2 停用
	RunStatus int    `json:"runStatus"` // 执行状态 1 开始执行 2 执行结束 3 执行异常 4 运行中
	ErrorCode int    `json:"errorCode"` // 错误码 0 正常 1 错误
	Msg       string `json:"msg"`       // 处理信息
	ProjName  string `json:"projName"`  // 工程名称

	ReTry   int    `json:"reTry"` // 重试次数
	ReCount int    `json:"reCount"`
	ChainId string `json:"chainId"` // 当前处理链路id
}

func (c *CronJob) GetKey() string {

	prefix := strings.HasPrefix(c.ID, "~KEY:")
	if prefix {
		return c.ID[len("~KEY:"):]
	}

	return c.Subject + ":" + c.FuncName + ":" + c.ID
}

func (c *CronJob) GetPrefix() string {
	return "cron_job"
}

func (c *CronJob) ToByte() []byte {
	return gjson.New(c).MustToJson()
}

func (c *CronJob) LoadBytes(bytes []byte) error {
	return gjson.New(bytes).Scan(c)
}
