package appllo_cfg

import (
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gcfg"
)

var _ gcfg.Adapter = (*Client)(nil)

var registerChangeAdapter = gmap.NewKVMap[string, ChangeAdapter](true)

func RegisterChangeAdapter(key string, adapter ChangeAdapter) {
	registerChangeAdapter.Set(key, adapter)
}

type ChangeAdapter interface {
	CheckKey(key string) bool // 是否进行更新
	OnChange(event *storage.ChangeEvent)
}

// OnChange is called when config changes.
func (c *Client) OnChange(event *storage.ChangeEvent) {
	for s, change := range event.Changes {
		c.value.Set(s, change.NewValue)
		registerChangeAdapter.Iterator(func(key string, value ChangeAdapter) bool {
			if value.CheckKey(s) {
				value.OnChange(event)
			}
			return true
		})
	}

	//	_ = c.updateLocalValue(gctx.New())
}
