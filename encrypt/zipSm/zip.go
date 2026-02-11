// 主要用于大文本的加密
package zipSm

import (
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/util/grand"
)

const (
	lkey = "myyixinfixkdsxdo"
	liv  = "rfgtecpoxserfghu"
)

var (
	localKey = lkey
	localIv  = liv
	bl       = 4
)

func SetKey(key string) {
	if len(key) != len(lkey) {
		panic("设置密钥长度不对")
		return
	}
	localKey = key
}
func SetIv(iv string) {
	if len(iv) != len(lkey) {
		panic("设置密钥长度不对")
		return
	}
	localIv = iv
}

func Encrypt(data []byte) (string, error) {

	gzip, err := gcompress.Gzip(data, grand.N(1, 9))
	if err != nil {
		return "", err
	}
	encrypt := crypto.FromBytes(gzip).SetKey(localKey).SM4().CBC().SetIv(localIv).PKCS7Padding().Encrypt()
	if err = encrypt.Error(); err != nil {
		return "", err
	}
	return encrypt.ToBase64String(), nil

}

func Decrypt(content string) ([]byte, error) {
	decrypt := crypto.FromBase64String(content).SetKey(localKey).SM4().CBC().SetIv(localIv).PKCS7Padding().Decrypt()
	if err := decrypt.Error(); err != nil {
		return nil, err
	}
	return gcompress.UnGzip(decrypt.ToBytes())

}
