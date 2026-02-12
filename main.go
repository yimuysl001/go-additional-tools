package main

import (
	"go-additional-tools/econf"
	"go-additional-tools/ecron/server/listen"
)

func main() {
	econf.MustInitConf()

	listen.StartNatsAndCron()
}
