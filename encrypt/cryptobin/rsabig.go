package cryptobin

import (
	"bytes"
	"github.com/deatil/go-cryptobin/cryptobin/rsa"
)

// 大数据加密
func RsaBigDataEncrypt(plainText, publicKey []byte, ecb ...bool) (cipherText []byte, err error) {
	rsatool := rsa.FromPublicKey(publicKey)

	pub := rsatool.GetPublicKey()
	pubSize, plainTextSize := pub.Size(), len(plainText)

	offSet, once := 0, pubSize-11

	buffer := bytes.Buffer{}
	for offSet < plainTextSize {
		endIndex := offSet + once
		if endIndex > plainTextSize {
			endIndex = plainTextSize
		}
		rsa2 := rsatool.FromBytes(plainText[offSet:endIndex])
		if ecb != nil && len(ecb) > 0 && ecb[0] {
			rsa2.EncryptECB()
		} else {
			rsa2 = rsa2.Encrypt()
		}

		err := rsa2.Error()
		if err != nil {
			return nil, err
		}

		bytesOnce := rsa2.ToBytes()

		buffer.Write(bytesOnce)
		offSet = endIndex
	}

	cipherText = buffer.Bytes()
	return cipherText, nil
}

// 大数据解密
func RsaBigDataDecrypt(cipherText, privateKey []byte, ecb ...bool) (plainText []byte, err error) {
	rsatool := rsa.FromPrivateKey(privateKey)
	pri := rsatool.GetPrivateKey()

	priSize, cipherTextSize := pri.Size(), len(cipherText)
	var offSet = 0
	var buffer = bytes.Buffer{}

	for offSet < cipherTextSize {
		endIndex := offSet + priSize
		if endIndex > cipherTextSize {
			endIndex = cipherTextSize
		}

		rsa2 := rsatool.FromBytes(cipherText[offSet:endIndex])
		if ecb != nil && len(ecb) > 0 && ecb[0] {
			rsa2.DecryptECB()
		} else {
			rsa2 = rsa2.Decrypt()
		}

		err := rsa2.Error()
		if err != nil {
			return nil, err
		}

		bytesOnce := rsa2.ToBytes()

		buffer.Write(bytesOnce)
		offSet = endIndex
	}

	plainText = buffer.Bytes()
	return plainText, nil
}

// 大数据私钥加密
func RsaBigDataEncryptByPri(plainText, privateKey []byte, ecb ...bool) (cipherText []byte, err error) {
	rsatool := rsa.FromPrivateKey(privateKey)

	pri := rsatool.GetPrivateKey()

	priSize, cipherTextSize := pri.Size(), len(plainText)

	offSet := 0

	buffer := bytes.Buffer{}
	for offSet < cipherTextSize {
		endIndex := offSet + priSize - 11
		if endIndex > cipherTextSize {
			endIndex = cipherTextSize
		}

		rsa2 := rsatool.FromBytes(plainText[offSet:endIndex])
		if ecb != nil && len(ecb) > 0 && ecb[0] {
			rsa2.PrivateKeyEncryptECB()
		} else {
			rsa2 = rsa2.PrivateKeyEncrypt()
		}

		err := rsa2.Error()
		if err != nil {
			return nil, err
		}

		bytesOnce := rsa2.ToBytes()

		buffer.Write(bytesOnce)
		offSet = endIndex
	}

	cipherText = buffer.Bytes()
	return cipherText, nil
}

// 大数据公钥解密
func RsaBigDataDecryptByPub(cipherText, publicKey []byte, ecb ...bool) (plainText []byte, err error) {
	rsatool := rsa.FromPublicKey(publicKey)

	pub := rsatool.GetPublicKey()
	pubSize, plainTextSize := pub.Size(), len(cipherText)

	offSet := 0

	buffer := bytes.Buffer{}
	for offSet < plainTextSize {
		endIndex := offSet + pubSize
		if endIndex > plainTextSize {
			endIndex = plainTextSize
		}

		rsa2 := rsatool.FromBytes(cipherText[offSet:endIndex])
		if ecb != nil && len(ecb) > 0 && ecb[0] {
			rsa2.PublicKeyDecryptECB()
		} else {
			rsa2 = rsa2.PublicKeyDecrypt()
		}
		err := rsa2.Error()
		if err != nil {
			return nil, err
		}

		bytesOnce := rsa2.ToBytes()

		buffer.Write(bytesOnce)
		offSet = endIndex
	}

	plainText = buffer.Bytes()
	return plainText, nil
}
