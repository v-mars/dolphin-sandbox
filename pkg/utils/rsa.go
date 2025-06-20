package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

/*
str := "你好123abca"
encryptstr, _ := RSAEncryptString(str, "files/public.pem")
fmt.Println(encryptstr)

decrypt, _ := RSADecryptString(encryptstr, "files/private.pem")
fmt.Println(decrypt)
*/

// RSAEncrypt RSA加密字节数组，返回字节数组
func RSAEncrypt(originalBytes []byte, filename string) ([]byte, error) {
	// 1、读取公钥文件，解析出公钥对象
	publicKey, err := ReadParsePublicKey(filename)
	if err != nil {
		return nil, err
	}
	// 2、RSA加密，参数是随机数、公钥对象、需要加密的字节
	// PKCS#1 v1.5 padding
	// Reader是一个全局共享的密码安全的强大的伪随机生成器
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, originalBytes)
}

// RSADecrypt RSA解密字节数组，返回字节数组
func RSADecrypt(cipherBytes []byte, filename string) ([]byte, error) {
	// 1、读取私钥文件，解析出私钥对象
	privateKey, err := ReadParsePrivaterKey(filename)
	if err != nil {
		return nil, err
	}
	// 2、ras解密，参数是随机数、私钥对象、需要解密的字节
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherBytes)
}

// ReadParsePublicKey2 读取公钥文件，解析出公钥对象
func ReadParsePublicKey2(filename string) (*rsa.PublicKey, error) {
	// 1、读取公钥文件，获取公钥字节
	publicKeyBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 2、解码公钥字节，生成加密块对象
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	// 3、解析DER编码的公钥，生成公钥接口
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 4、公钥接口转型成公钥对象
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey, nil
}

func ReadParsePublicKey(pk string) (*rsa.PublicKey, error) {
	// 2、解码公钥字节，生成加密块对象
	block, _ := pem.Decode([]byte(pk))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	// 3、解析DER编码的公钥，生成公钥接口
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 4、公钥接口转型成公钥对象
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey, nil
}

// ReadParsePrivaterKey2 读取私钥文件，解析出私钥对象
func ReadParsePrivaterKey2(filename string) (*rsa.PrivateKey, error) {
	// 1、读取私钥文件，获取私钥字节
	privateKeyBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 2、对私钥文件进行编码，生成加密块对象
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("sk私钥信息错误！")
	}
	// 3、解析DER编码的私钥，生成私钥对象
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// ReadParsePrivaterKey 读取私钥文件，解析出私钥对象
func ReadParsePrivaterKey(sk string) (*rsa.PrivateKey, error) {
	// 2、对私钥文件进行编码，生成加密块对象
	block, _ := pem.Decode([]byte(sk))
	if block == nil {
		return nil, errors.New("sk私钥信息错误！")
	}
	// 3、解析DER编码的私钥，生成私钥对象
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// RSAEncryptString RSA加密字符串，返回base64处理的字符串
func RSAEncryptString(originalText, filename string) (string, error) {
	cipherBytes, err := RSAEncrypt([]byte(originalText), filename)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// RSADecryptString RSA 解密经过base64处理的加密字符串，返回加密前的明文
func RSADecryptString(cipherlText, filename string) (string, error) {
	cipherBytes, _ := base64.StdEncoding.DecodeString(cipherlText)
	originalBytes, err := RSADecrypt(cipherBytes, filename)
	if err != nil {
		return "", err
	}
	return string(originalBytes), nil
}
