package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

type LocalCipher struct {
	key   []byte
	ivAes []byte
}

func InitCipher(key []byte, ivAes []byte) (c *LocalCipher, err error) {
	c = new(LocalCipher)
	if len(key) != 32 {
		return nil, errors.New("key格式错误(32byte)")
	}
	if len(ivAes) != 16 {
		return nil, errors.New("ivAes格式错误(16byte)")
	}
	c.key = key
	c.ivAes = ivAes
	return
}

func (c *LocalCipher) Encrypt(plainText []byte) (string, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	padding := block.BlockSize() - (len(plainText) % block.BlockSize())
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	paddingText := append(plainText, padText...)

	blockMode := cipher.NewCBCEncrypter(block, c.ivAes)
	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)
	return base64.RawURLEncoding.EncodeToString(cipherText), nil
}

func (c *LocalCipher) Decrypt(cipherText string) (string, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	tmp, err := base64.RawURLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", errors.New("base64解码失败")
	}

	blockMode := cipher.NewCBCDecrypter(block, c.ivAes)
	paddingText := make([]byte, len(tmp))
	blockMode.CryptBlocks(paddingText, tmp)

	length := len(paddingText)
	number := int(paddingText[length-1])
	if number > length {
		return "", errors.New("ErrPaddingSize错误")
	}

	return string(paddingText[:length-number]), nil
}

func Sha256Hex(data []byte) string {
	digest := sha256.New()
	digest.Write(data)
	return hex.EncodeToString(digest.Sum(nil))
}
