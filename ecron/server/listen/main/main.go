package main

import (
	"go-additional-tools/econf"
	"go-additional-tools/ecron/server/listen"
)

func main() {
	econf.MustInitConf("E:\\wk\\git\\go-additional-tools\\ecron\\server\\listen\\main\\config.conf")

	listen.StartNatsAndCron()
}
