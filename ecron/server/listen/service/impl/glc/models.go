package sysmnt

import "github.com/syndtr/goleveldb/leveldb"

type StoreModel interface {
	GetKey() string
	GetPrefix() string
	ToByte() []byte
	LoadBytes(bytes []byte) error
}
type SysmntStorage struct {
	storeName string
	subPath   string      // 存储目录下的相对路径（存放数据）
	leveldb   *leveldb.DB // leveldb
	lastTime  int64       // 最后一次访问时间
	closing   bool        // 是否关闭中状态
}
