package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptionService Push Token加密服务
type EncryptionService struct {
	key []byte // 32字节密钥（AES-256）
}

// NewEncryptionService 创建加密服务
func NewEncryptionService(keyStr string) (*EncryptionService, error) {
	var key []byte

	// 尝试base64解码
	decoded, err := base64.StdEncoding.DecodeString(keyStr)
	if err == nil && len(decoded) == 32 {
		// 成功解码且长度为32字节，使用解码后的密钥
		key = decoded
	} else if len(keyStr) == 32 {
		// 直接是32字节的原始字符串
		key = []byte(keyStr)
	} else {
		return nil, fmt.Errorf("encryption key must be 32 bytes (raw string or base64 encoded)")
	}

	return &EncryptionService{
		key: key,
	}, nil
}

// Encrypt 加密Push Token
func (s *EncryptionService) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密Push Token
func (s *EncryptionService) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("invalid ciphertext")
	}

	nonce := data[:gcm.NonceSize()]
	ciphertext = string(data[gcm.NonceSize():])

	plaintext, err := gcm.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
