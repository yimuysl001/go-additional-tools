package sysmnt

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/syndtr/goleveldb/leveldb"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	mapStorage   sync.Map
	mapStorageMu sync.Mutex
	sdbMu        sync.Mutex
)

const (
	MODEL_PREFIX = "sysmodel:"
	MODEL_KEYS   = "syskeys:"
)

func NewSysmntStorage(ctx context.Context, storeName string) *SysmntStorage {
	mapStorageMu.Lock()         // 缓存map锁
	defer mapStorageMu.Unlock() // 缓存map解锁
	subPath := "sysmnt"

	cacheName := path.Join(storeName, subPath)
	cacheStore := getCacheStore(cacheName)
	if cacheStore != nil && !cacheStore.IsClose() {
		return cacheStore
	}
	// 缓存有则取用
	store := new(SysmntStorage)
	store.subPath = cacheName
	store.closing = false
	store.storeName = storeName

	store.lastTime = time.Now().Unix()

	dbPath := g.Cfg().MustGet(ctx, "sysmnt.subpath", "storage")

	db, err := leveldb.OpenFile(path.Join(dbPath.String(), store.subPath), nil) // 打开（在指定子目录中存放数据）
	if err != nil {
		g.Log().Error(ctx, "打开leveldb失败："+err.Error())
		panic(err)
	}
	store.leveldb = db

	mapStorage.Store(cacheName, store) // 缓存起来
	// 逐秒判断，若闲置超时则自动关闭
	go store.autoCloseWhenMaxIdle()
	g.Log().Info(ctx, "打开SysmntStorage："+store.subPath)
	return store
}

func getCacheStore(cacheName string) *SysmntStorage {
	value, ok := mapStorage.Load(cacheName)
	if ok && value != nil {
		cacheStore := value.(*SysmntStorage)
		if !cacheStore.IsClose() {
			return cacheStore // 缓存中未关闭的存储对象
		}
	}
	return nil
}

func (s *SysmntStorage) Close() {
	if s == nil || s.closing { // 优雅退出时可能会正好nil，判断一下优雅点
		return
	}

	sdbMu.Lock()         // 锁
	defer sdbMu.Unlock() // 解锁
	if s.closing {
		return
	}

	s.closing = true
	s.leveldb.Close()

	mapStorage.Delete(s.storeName)
	g.Log().Info(gctx.GetInitCtx(), "关闭SysmntStorage："+s.subPath)

}

// 是否关闭中状态
func (s *SysmntStorage) IsClose() bool {
	return s.closing
}

// 是否关闭中状态
func (s *SysmntStorage) GetAllKeys(ctx context.Context, prefix string) []string {
	bs, err := s.Get([]byte(MODEL_KEYS + prefix))
	if err != nil {
		var rs []string
		return rs
	}
	return strings.Split(string(bs), ",")
}

func (s *SysmntStorage) SaveAllKeys(ctx context.Context, prefix string, names []string) error {
	return s.Put([]byte(MODEL_KEYS+prefix), []byte(strings.Join(names, ",")))
}

func (s *SysmntStorage) GetModel(ctx context.Context, p StoreModel) (err error, bool2 bool) {
	bs, err := s.Get([]byte(MODEL_PREFIX + p.GetPrefix() + ":" + p.GetKey()))
	if err != nil {
		return err, false
	}
	return p.LoadBytes(bs), true
}

func (s *SysmntStorage) SaveModel(ctx context.Context, p StoreModel) error {
	toByte := p.ToByte()
	_, b := s.GetModel(ctx, p)
	if !b { // 新增
		keys := s.GetAllKeys(ctx, p.GetPrefix())
		keys = append(keys, p.GetKey())
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j] // 排序
		})
		s.SaveAllKeys(ctx, p.GetPrefix(), keys)

	}
	return s.Put([]byte(MODEL_PREFIX+p.GetPrefix()+":"+p.GetKey()), toByte)
}

func (s *SysmntStorage) DelModel(ctx context.Context, p StoreModel) error {
	var newNames []string
	names := s.GetAllKeys(ctx, p.GetPrefix())

	for i, name := range names {
		if name == p.GetKey() {
			if i != len(names)-1 {
				newNames = append(newNames, names[i+1:]...)
			}
			break
		} else {
			newNames = append(newNames, name)
		}
	}
	s.SaveAllKeys(ctx, p.GetPrefix(), newNames)
	return s.Del([]byte(MODEL_PREFIX + p.GetPrefix() + ":" + p.GetKey()))

}

func stringToBytes(str string) []byte {
	return []byte(str)
}

func (s *SysmntStorage) GetStorageDataCount(storeName string) uint32 {
	bt, err := s.Get(stringToBytes("data:" + storeName))
	if err != nil {
		return 0
	}
	return gconv.Uint32(bt)
}

func (s *SysmntStorage) SetStorageDataCount(storeName string, count uint32) {
	s.Put(stringToBytes("data:"+storeName), stringToBytes(gconv.String(count)))
}

func (s *SysmntStorage) GetStorageIndexCount(storeName string) uint32 {
	bt, err := s.Get(stringToBytes("index:" + storeName))
	if err != nil {
		return 0
	}
	return gconv.Uint32(bt)
}

func (s *SysmntStorage) SetStorageIndexCount(storeName string, count uint32) {
	s.Put(stringToBytes("index:"+storeName), stringToBytes(gconv.String(count)))
}

func (s *SysmntStorage) DeleteStorageInfo(storeName string) error {
	err := s.leveldb.Delete(stringToBytes("data:"+storeName), nil)
	if err != nil {
		return err
	}
	err = s.leveldb.Delete(stringToBytes("index:"+storeName), nil)
	if err != nil {
		return err
	}
	return nil
}

// 直接存入数据到leveldb
func (s *SysmntStorage) Put(key []byte, value []byte) error {
	if s.closing {
		return errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Put(key, value, nil)
}

// 直接从leveldb取数据
func (s *SysmntStorage) Get(key []byte) ([]byte, error) {
	if s.closing {
		return nil, errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Get(key, nil)
}

// 直接从leveldb取数据
func (s *SysmntStorage) Del(key []byte) error {
	if s.closing {
		return errors.New("current storage is closed") // 关闭中或已关闭时拒绝服务
	}
	s.lastTime = time.Now().Unix()
	return s.leveldb.Delete(key, nil)
}
func (s *SysmntStorage) autoCloseWhenMaxIdle() {
	var maxidle = g.Cfg().MustGet(gctx.GetInitCtx(), "sysmnt.maxidle", "").Duration()

	if maxidle > 0 {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			if time.Now().Unix()-s.lastTime > int64(maxidle) {
				s.Close()
				ticker.Stop()
				break
			}
		}
	}
}
