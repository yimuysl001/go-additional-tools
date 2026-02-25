package egoja

import (
	"bytes"
	"compress/flate"
	"context"
	rsae "crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/dop251/goja"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/pierrec/lz4/v4"
	"go-additional-tools/egoja/pkgs"
	"go-additional-tools/egoja/require"
	"io"
	"sync"
	"testing"
)

var scripts = map[string]string{
	"mathUtils": `
            export function add(a, b) { return a + b; }
            export function multiply(a, b) { return a * b; }
        `,
	"dataProcessor": `
            import { add, multiply } from 'mathUtils.js';
            export function process(data) {
                return data.map(item => multiply(item.value, add(item.base, 1)));
            };
			export function processOne(data) {
                return multiply(data.value, add(data.base, 1));
            };
        `,
}

func TestScript(t *testing.T) {

	for name, source := range scripts {
		require.RegisterFuncScript(name, source)
	}
	data := make([]map[string]any, 100)
	for i := 0; i < 100; i++ {
		data[i] = map[string]any{"value": i + 1, "base": i + 2}
	}
	var wg sync.WaitGroup
	for i, datum := range data {
		wg.Add(1)
		item := datum
		go func() {
			defer wg.Done()

			id, err := ExecScriptById(gctx.New(), "test", `import {processOne} from 'dataProcessor.js'
return processOne(data)`, map[string]any{"data": item})
			fmt.Println(i, "结果：", id, err)
		}()

	}
	wg.Wait()
}

func TestScript2(t *testing.T) {

	for name, source := range scripts {
		require.RegisterFuncScript(name, source)
	}
	data := make([]map[string]any, 100)
	for i := 0; i < 100; i++ {
		data[i] = map[string]any{"value": i + 1, "base": i + 2}
	}
	var wg sync.WaitGroup
	for i, datum := range data {
		wg.Add(1)
		item := datum
		go func() {
			defer wg.Done()
			id, err := ExecScript(gctx.New(), `import {processOne} from 'dataProcessor.js'
return processOne(data)`, map[string]any{"data": item})
			fmt.Println(i, "结果：", id, err)
		}()

	}
	wg.Wait()
}

func TestScript3(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := ExecScriptById(gctx.New(), "test", `
return index+1
`, map[string]any{"index": i})
			fmt.Println("结果：", id, err)
		}()
	}
	wg.Wait()
}

func TestScript4(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := ExecScriptById(gctx.New(), "test", `
var ccc =  index+1
return  ctx.Value("vm").Get("ccc")
`, map[string]any{"index": i})
			fmt.Println("结果：", id, err)
		}()
	}
	wg.Wait()
}

func TestScript5(t *testing.T) {
	//	id, err := ExecScriptById(gctx.New(), "test", `
	//  var foo = "bar";
	//return  ctx.Value("vm").Get("foo")
	//`, map[string]any{"index": 123})
	//	fmt.Println("结果：", id, err)

	vm := goja.New()

	// 注册一个函数，通过闭包持有 vm
	vm.Set("getGlobal", func(call goja.FunctionCall) goja.Value {
		name := call.Arguments[0].String()
		val := vm.Get(name) // 使用闭包捕获的 vm
		return val
	})

	// 脚本中调用
	a, e := vm.RunString(`(function() {
 				  var foo = "bar";
    return getGlobal("foo"); // 输出 "bar"
        })()
  
`)
	fmt.Println(a, e)

}

func TestGoScript(t *testing.T) {
	pkgs.GoEnv()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := ExecScriptById(gctx.New(), "testgo", `
import * as strings from 'strings'
var sb1 = strings.Builder
var sb2 = strings.Builder
sb1.WriteString("123==>")

for (var i = 0; i < 10; i++) { 
  sb1.WriteString((index+i)+'')
 }

sb2.WriteString("=================\n")
sb2.WriteString(sb1.String())
return sb2.String()
`, map[string]any{"index": i})
			fmt.Println("结果：", id, err)
		}()
	}
	wg.Wait()
}

func BenchmarkVariableP(b *testing.B) {
	require.RegisterCommonParameter("targs", func(b string, args ...any) {
		fmt.Println(b)
		fmt.Println(args...)
	})

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := ExecScriptById(context.TODO(), "AAA", `var aaa= ["a","b","ccc",ddd]
targs("a",  ...aaa)
`, map[string]any{
				"ddd": "123546",
			})
			fmt.Println(err)
		}()

	}
	wg.Wait()

}

const (
	rsaprivkey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDVwXbbYTfhe9X1
L9jvfDQkoxuXs5ef5Vp21BnudaCwFvNPMphgSyavgySujm0uzTi7rKkXIJpS1WyP
Z8SKZ/nKJgnOZ80qm1ZeCVefXD48N6xCatYUbzc53jibUT9xhdxZDZC39d1olGcI
6y9tvmpD/7skRJnPnOEUNqm6MzfR9u69RUgU8nGoXuiKi/rsOTq/9hBrlO3y1xhU
FyJbHx/Jn3ibe2JRB4v/Arsss1Swh8hk6QrPbYdjI9emOdjDJxdxY4/8VUKnXkdL
XJ8ljcnjg54bQulQdRaTlJg7mhs3XEOI2+pR05R84XSST0S5+FSROgjKXMeJ1HOb
AB8A7gbLAgMBAAECggEAHGZd56aHtCUDN6upw6KEPbHmF/ERM3oCMPOGlnt5DqxI
m2ck06ETJ2fh1xyuk59J6wrKVHuicAr8fM5+pcehanhEqUXP/oSGTxbSq3un+9GH
HjPvnuYVN6e3bFC5Ky4WfpvMydmYFy2lxYuUe5zm+tiJVreT7/+Z8AOV6AVSD0qP
ZKn5fq6DmLcntPIQsb8OGGFmq/Yig21ShbTxa7JgexOfNZTs6qZvRFodObPOwXlj
r8xit4g5c4J71ZMW1+amwmvYAccFuXF5TlVHw+xMkbm3ySMHsiS5MMlCr1/ECmPb
xJjXQHMBS8Tfrc9mj6epRWW9kAxAAES3EDY+sVUjoQKBgQDegKZOvGegMI9LpLty
yf2F15WGIldIklH0dSkUrtTTFprE7/uhK8r2QwO43m8odGN5YqgSHnIHx9isq7t1
PU4kOcBP3RzZcsgG0O+0MpCOAgovJjZ2krHXvHcZWm6hCjgMu/lZPhJgHrq2lalo
Yz75V7NKl03GEa1WkoLn9qM4WQKBgQD177VQUBgg5EHUy39+tiUNp7Dy83cyv171
zKOO/y+eSr3WlPlrMNvwJa4DY56xesK5PzVuadM65UlsUwVYKiwEN+7+QaC9FMeQ
WvMeTk4mM7qAVlIaWTNcABsZrxjp0LWQxad55ChUVTBh01uBdc27qpwyW3Xd6y+K
pnLv+1eTwwKBgDDqh6V3tjB5fIdcx/kMfzgVlUHP+vBxeqMLvuRVK2Tc61mwiNl+
DzjkssTJ4hY6wEPHdLvHBbrALNqJRsUXnT5JlAX6zoTfvyoAdTJgi3cs66BB/mdD
COYtAOIKB5hP7tKd4MvF4bRQDSxm6r+QUh/vL/OOIAMTj9Aglbb5ehjBAoGBALlz
j4rHStqKpOWcqkBXg2tflzwswSagTjAVpwQug67ed3Z3Efl1d3QIRcbCeSkmA+4C
rvzaifDwc0Re+jm4W0a3Et3hiR7rq2y8WHXy4FVITot2DCVYPDVU0xq0AZpWyoMn
uJlepdaqAnjSEz91ILUx+uSyORglv8zSpPs30ZtXAoGARP+3TLxyVxzoLLvSdCda
hKjvYwHgf23fdp8MM8XcAWQRJ0+o+Ts369X8pujKsoSw9VgjKiVtcgJ3bZf1CdRm
gzWOryiHQK5txwEog3TqFko/UVWh2Gqg//okxlTGj9/nkiA6FjQLzPPjlrJ2/TV8
mpvglZ+GwFhybnUTqRUPFyg=
-----END PRIVATE KEY-----`
	respubkey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1cF222E34XvV9S/Y73w0
JKMbl7OXn+VadtQZ7nWgsBbzTzKYYEsmr4Mkro5tLs04u6ypFyCaUtVsj2fEimf5
yiYJzmfNKptWXglXn1w+PDesQmrWFG83Od44m1E/cYXcWQ2Qt/XdaJRnCOsvbb5q
Q/+7JESZz5zhFDapujM30fbuvUVIFPJxqF7oiov67Dk6v/YQa5Tt8tcYVBciWx8f
yZ94m3tiUQeL/wK7LLNUsIfIZOkKz22HYyPXpjnYwycXcWOP/FVCp15HS1yfJY3J
44OeG0LpUHUWk5SYO5obN1xDiNvqUdOUfOF0kk9EufhUkToIylzHidRzmwAfAO4G
ywIDAQAB
-----END PUBLIC KEY-----`
)

func TestRsa(t *testing.T) {

	//key := rsa.GenerateKey(2048)
	//
	//fmt.Println(getPrivateKeyPKCS8String(key.GetPrivateKey()))
	//
	//fmt.Println(getPublicKeyString(key.GetPublicKey()))
	data := []byte("test-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-pass3333333333333333333333333333333333333333333333333333test-pa2222222222222222222222222222222222222222222sstest-passt111111111111111111111111111111111est-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-passtest-pass")

	text, err := pkgs.RsaBigDataEncrypt(data, []byte(respubkey))

	fmt.Println(base64.StdEncoding.EncodeToString(text), err)
	gzip, err := FlateCompress(text)

	var ziphex = base64.StdEncoding.EncodeToString(gzip)
	fmt.Println(ziphex, err)
	decodeString, _ := base64.StdEncoding.DecodeString(ziphex)
	unGzip, err := FlateDecompress(decodeString)
	plainText, err := pkgs.RsaBigDataDecrypt(unGzip, []byte(rsaprivkey))

	fmt.Println(string(plainText), err)

}

// 获取私钥字符串（PKCS8格式）
func getPrivateKeyPKCS8String(privateKey *rsae.PrivateKey) (string, error) {
	// 将私钥转换为PKCS8格式的DER编码
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	// 创建PEM块
	privateKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// 将PEM块编码为字符串
	return string(pem.EncodeToMemory(privateKeyPEM)), nil
}

// 获取公钥字符串
func getPublicKeyString(publicKey *rsae.PublicKey) (string, error) {
	// 将公钥转换为PKIX格式的DER编码
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	// 创建PEM块
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	// 将PEM块编码为字符串
	return string(pem.EncodeToMemory(publicKeyPEM)), nil
}

func FlateCompress(s []byte) ([]byte, error) {
	var buf bytes.Buffer
	fw, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = fw.Write(s)
	if err != nil {
		return nil, err
	}

	if err := fw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func FlateDecompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	fr := flate.NewReader(buf)
	defer fr.Close()

	decompressed, err := io.ReadAll(fr)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

func LZ4Compress(s []byte) ([]byte, error) {
	// 预分配足够空间
	buf := make([]byte, lz4.CompressBlockBound(len(s)))
	n, err := lz4.CompressBlock(s, buf, nil)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func LZ4Decompress(compressed []byte, originalSize int) ([]byte, error) {
	dst := make([]byte, originalSize)
	_, err := lz4.UncompressBlock(compressed, dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
