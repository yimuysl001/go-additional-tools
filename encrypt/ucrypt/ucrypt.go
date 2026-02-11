package ucrypt

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
)

const (
	C1       = 28853
	C2       = 31836
	PassKey  = 5728
	RightKey = 35762
)

var (
	base64Map [256]byte
	encodeMap = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
)

func init() {
	for i := range base64Map {
		base64Map[i] = 0
	}
	base64Map['+'] = 62
	base64Map['/'] = 63
	for i := 0; i < 10; i++ {
		base64Map['0'+i] = byte(52 + i)
	}
	for i := 0; i < 26; i++ {
		base64Map['A'+i] = byte(i)
	}
	for i := 0; i < 26; i++ {
		base64Map['a'+i] = byte(26 + i)
	}

	//fmt.Println(base64Map)
}

func decodeBase64(S []byte) []byte {
	var result []byte
	switch len(S) {
	case 2:
		var (
			f1 = uint32(base64Map[S[0]])
			f2 = uint32(base64Map[S[1]]) << 6
		)
		i := f1 + f2
		result = []byte{byte(i & 0xFF)}
	case 3:
		var (
			f1 = uint32(base64Map[S[0]])
			f2 = uint32(base64Map[S[1]]) << 6
			f3 = uint32(base64Map[S[2]]) << 12
		)
		i := f1 + f2 + f3
		result = []byte{
			byte(i & 0xFF),
			byte((i >> 8) & 0xFF),
		}
	case 4:
		var (
			f1 = uint32(base64Map[S[0]])
			f2 = uint32(base64Map[S[1]]) << 6
			f3 = uint32(base64Map[S[2]]) << 12
			f4 = uint32(base64Map[S[3]]) << 18
		)
		i := f1 + f2 + f3 + f4
		result = []byte{
			byte(i & 0xFF),
			byte((i >> 8) & 0xFF),
			byte((i >> 16) & 0xFF),
		}
	default:
		return nil
	}
	return result
}

func encodeBase64(S []byte) []byte {
	var i uint32
	switch len(S) {
	case 1:
		i = uint32(S[0])
		return []byte{
			encodeMap[i%64],
			encodeMap[(i>>6)%64],
		}
	case 2:
		i = uint32(S[0]) | uint32(S[1])<<8
		return []byte{
			encodeMap[i%64],
			encodeMap[(i>>6)%64],
			encodeMap[(i>>12)%64],
		}
	case 3:
		i = uint32(S[0]) | uint32(S[1])<<8 | uint32(S[2])<<16
		return []byte{
			encodeMap[i%64],
			encodeMap[(i>>6)%64],
			encodeMap[(i>>12)%64],
			encodeMap[(i>>18)%64],
		}
	default:
		return nil
	}
}

func PostProcess(src []byte) string {
	var result []byte
	for i := 0; i < len(src); {
		chunkSize := 3
		if len(src)-i < chunkSize {
			chunkSize = len(src) - i
		}
		chunk := src[i : i+chunkSize]
		encoded := encodeBase64(chunk)
		result = append(result, encoded...)
		i += chunkSize
	}
	return string(result)
}

func InternalEncrypt(s []byte, key int) []byte {
	seed := int64(key)
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		encrypted := byte((int64(s[i]) ^ (seed >> 8)) & 0xFF)
		result[i] = encrypted
		seed = (seed+int64(encrypted))*C1 + C2
	}
	return result
}

func PreProcess(ss []byte) []byte {
	var result []byte
	for i := 0; i < len(ss); {
		chunkSize := 4
		if len(ss)-i < chunkSize {
			chunkSize = len(ss) - i
		}
		chunk := ss[i : i+chunkSize]
		if i < len(ss) {
			chunk[0] = ss[i] & 0xFF
		}
		i++
		if i < len(ss) {
			chunk[1] = ss[i] & 0xFF
		}
		i++
		if i < len(ss) {
			chunk[2] = ss[i] & 0xFF
		}
		i++
		if i < len(ss) {
			chunk[3] = ss[i] & 0xFF
		}
		i++
		decoded := decodeBase64(chunk)
		result = append(result, decoded...)
	}

	//for i := 0; i < len(ss); {
	//	chunkSize := 4
	//	if len(ss)-i < chunkSize {
	//		chunkSize = len(ss) - i
	//	}
	//	chunk := ss[i : i+chunkSize]
	//	decoded := decodeBase64(chunk)
	//	result = append(result, decoded...)
	//	i += chunkSize
	//}
	return result
}

func InternalDecrypt(s []byte, key int) []byte {
	seed := int64(key)
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		decrypted := byte((int64(s[i]) ^ (seed >> 8)) & 0xFF)
		result[i] = decrypted
		seed = (seed+int64(s[i]))*C1 + C2
	}
	return result
}

func Decrypt(s string, key int) string {
	decrypted := InternalDecrypt([]byte(s), key)
	preProcessed := PreProcess(decrypted)
	return string(preProcessed)
}

func Encrypt(s string, key int) string {
	encrypted := InternalEncrypt([]byte(s), key)
	return PostProcess(encrypted)
}

func Decrypt2(inStr, keyStr string, key int) (string, error) {
	keyBytes, err := toGBK(keyStr)
	if err != nil {
		return "", err
	}
	A := key
	for _, b := range keyBytes {
		A += int(b)
	}

	inBytes := []byte(inStr)
	preProcessed := PreProcess(inBytes)
	decrypted := InternalDecrypt(preProcessed, A)

	var cBytes []byte
	for _, b := range decrypted {
		cBytes = append(cBytes, b^0xA5)
	}

	result, err := fromGBK(cBytes)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func Encrypt2(inStr, keyStr string, key int) (string, error) {
	keyBytes, err := toGBK(keyStr)
	if err != nil {
		return "", err
	}
	A := key
	for _, b := range keyBytes {
		A += int(b)
	}

	inBytes, err := toGBK(inStr)
	if err != nil {
		return "", err
	}

	var cBytes []byte
	for _, b := range inBytes {
		cBytes = append(cBytes, b^0xA5)
	}

	encrypted := InternalEncrypt(cBytes, A)
	result := PostProcess(encrypted)

	return result, nil
}

func toGBK(s string) ([]byte, error) {
	encoder := simplifiedchinese.GBK.NewEncoder()
	return io.ReadAll(transform.NewReader(bytes.NewReader([]byte(s)), encoder))
}

func fromGBK(b []byte) ([]byte, error) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	return io.ReadAll(transform.NewReader(bytes.NewReader(b), decoder))
}
