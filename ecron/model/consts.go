package model

const (
	CronNats              = "cron-subs"
	ExecNats              = "cron-exec"
	CronCheckLockNatsSubj = "cron.checkLock"
	CronCheckIsLock       = "cron.isLock"
	CronCheckNoLock       = "cron.isNoLock"

	CronAddNatsSubj    = "cron.add"
	CronDelNatsSubj    = "cron.del"
	CronListNatsSubj   = "cron.list"
	CronUpDateNatsSubj = "cron.update"
	CronStopNatsSubj   = "cron.stop"
	CronStartNatsSubj  = "cron.start"
)

type CronSubj string

const (
	CronAddNatsSubjType    = CronSubj(CronAddNatsSubj)
	CronDelNatsSubjType    = CronSubj(CronDelNatsSubj)
	CronUpDateNatsSubjType = CronSubj(CronUpDateNatsSubj)
	CronStopNatsSubjType   = CronSubj(CronStopNatsSubj)
	CronStartNatsSubjType  = CronSubj(CronStartNatsSubj)
)
