package econf

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
	"time"
)

func TestCfg(t *testing.T) {
	MustInitConf()

	get, err := g.Cfg().Get(gctx.New(), "pid_file")

	fmt.Println(get, err)

	get, err = g.Cfg().Get(gctx.New(), "cron.port")

	fmt.Println(get, err)

	get, err = g.Cfg().Get(gctx.New(), "sql-query")

	fmt.Println(get, err)

	get, err = g.Cfg().Get(gctx.New(), "test")

	fmt.Println(get, err)

	for {
		get, err = g.Cfg().Get(gctx.New(), "cron.port")

		fmt.Println(get, err)
		time.Sleep(5 * time.Second)

	}

}
