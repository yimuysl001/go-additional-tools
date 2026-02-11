package econf

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
)

func TestCfg(t *testing.T) {
	MustInitConf()

	get, err := g.Cfg().Get(gctx.New(), "pid_file")

	fmt.Println(get, err)

	get, err = g.Cfg().Get(gctx.New(), "cron.port")

	fmt.Println(get, err)

	get, err = g.Cfg().Get(gctx.New(), "sql-query")

	fmt.Println(get, err)

}
