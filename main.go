package main

import (
	"go-additional-tools/econf"
	"go-additional-tools/ecron/server/listen"
)

// //go:generate goversioninfo -icon=icon.ico
func main() {
	econf.MustInitConf()

	listen.StartNatsAndCron()
}
