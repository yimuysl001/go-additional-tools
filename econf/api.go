package econf

import (
	"github.com/gogf/gf/v2/frame/g"
)

const (
	// 默认加密密钥和IV
	defaultKey = "myyixinfixkdsxdo"
	defaultIv  = "rfgtecpoxserfghu"

	// 加密相关常量
	encryptionPrefix = "ENC("
	encryptionSuffix = ")"
	encryptionKey    = 5728

	// 配置文件相关常量
	defaultConfigName = "config"
	configFilePathKey = "gf.gcfg.file"
)

var (
	localKey = ""
	localIv  = ""
)

// ConfigOption 配置选项
type ConfigOption struct {
	Key  string
	Iv   string
	Name string
}

// WithKey 设置自定义密钥
func WithKey(key string) ConfigOption {
	return ConfigOption{Key: key}
}

func WithName(name string) ConfigOption {
	return ConfigOption{Name: name}
}

// WithIv 设置自定义IV
func WithIv(iv string) ConfigOption {
	return ConfigOption{Iv: iv}
}

// SetLocalKey 设置本地密钥（已废弃，建议使用WithKey）
// Deprecated: 使用 InitConf(file, WithKey(key)) 替代
func SetLocalKey(ckey string) {
	localKey = ckey
}

// SetLocalVi 设置本地IV（已废弃，建议使用WithIv）
// Deprecated: 使用 InitConf(file, WithIv(iv)) 替代
func SetLocalVi(civ string) {
	localIv = civ
}

// getLocalKey 获取本地密钥
func getLocalKey() string {
	if localKey == "" {
		return defaultKey
	}
	return localKey
}

// getLocalIv 获取本地IV
func getLocalIv() string {
	if localIv == "" {
		return defaultIv
	}
	return localIv
}

// InitConf 初始化配置，注意需最优先加载
// 支持可选参数：WithKey(), WithIv()
func InitConf(files []string, opts ...ConfigOption) error {
	cfg := g.Cfg()
	// 应用配置选项
	for _, opt := range opts {
		if opt.Key != "" {
			localKey = opt.Key
		}
		if opt.Iv != "" {
			localIv = opt.Iv
		}
		if opt.Name != "" {
			cfg = g.Cfg(opt.Name)
		}

	}

	adapter := NewAdapter(files...)
	cfg.SetAdapter(adapter)
	return nil
}

// MustInitConf 初始化配置，失败时panic
func MustInitConf(files ...string) {
	if err := InitConf(files); err != nil {
		panic(err)
	}
}
