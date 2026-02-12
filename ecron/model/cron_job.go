package model

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"strings"
)

type CronJob struct {
	ID          string `json:"id"`          // 主键ID
	FuncName    string `json:"func_name"`   // 执行函数名
	Description string `json:"description"` // 任务描述
	CronExpr    string `json:"cron_expr"`   // Cron表达式
	Subject     string `json:"subject"`     // NATS主题/分组
	Params      string `json:"params"`      // 执行参数
	Timeout     string `json:"timeout"`     // 超时时间
	CheckLock   bool   `json:"check_lock"`  // 锁检查开关

	LastRun   string `json:"last_run"`   // 最后执行时间
	LastEnd   string `json:"last_end"`   // 执行完成时间
	Status    int    `json:"status"`     // 任务状态 1 启用 2 停用
	RunStatus int    `json:"run_status"` // 执行状态 1 开始执行 2 执行结束 3 执行异常 4 运行中
	ErrorCode int    `json:"error_code"` // 错误码 0 正常 1 错误
	Msg       string `json:"msg"`        // 处理信息
	ProjName  string `json:"proj_name"`  // 工程名称

	ReTry   int    `json:"re_try"` // 重试次数
	ReCount int    `json:"re_count"`
	ChainId string `json:"chain_id"` // 当前处理链路id
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
