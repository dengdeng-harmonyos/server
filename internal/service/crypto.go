package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
)

// EncryptedMessage 加密后的消息结构
type EncryptedMessage struct {
	EncryptedAESKey  string `json:"encrypted_aes_key"` // RSA加密的AES密钥
	EncryptedContent string `json:"encrypted_content"` // AES加密的消息内容
	IV               string `json:"iv"`                // AES IV向量
}

// MessageContent 原始消息内容
type MessageContent struct {
	Title   string                 `json:"title"`
	Content string                 `json:"content"`
	Data    map[string]interface{} `json:"data"`
}

// CryptoService 加密服务
type CryptoService struct{}

// NewCryptoService 创建加密服务实例
func NewCryptoService() *CryptoService {
	return &CryptoService{}
}

// EncryptMessage 使用RSA+AES混合加密消息
// publicKeyPEM: PEM格式的RSA公钥
// message: 要加密的消息内容
func (s *CryptoService) EncryptMessage(publicKeyPEM string, message MessageContent) (*EncryptedMessage, error) {
	// 1. 将消息序列化为JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return nil, errors.New("failed to marshal message: " + err.Error())
	}

	// 2. 生成随机AES密钥（32字节 = AES-256）
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, errors.New("failed to generate AES key: " + err.Error())
	}

	// 3. 使用AES-GCM加密消息内容
	encryptedContent, iv, err := s.encryptWithAES(aesKey, messageJSON)
	if err != nil {
		return nil, err
	}

	// 4. 使用RSA公钥加密AES密钥
	encryptedAESKey, err := s.encryptWithRSA(publicKeyPEM, aesKey)
	if err != nil {
		return nil, err
	}

	return &EncryptedMessage{
		EncryptedAESKey:  base64.StdEncoding.EncodeToString(encryptedAESKey),
		EncryptedContent: base64.StdEncoding.EncodeToString(encryptedContent),
		IV:               base64.StdEncoding.EncodeToString(iv),
	}, nil
}

// encryptWithAES 使用AES-GCM加密数据
func (s *CryptoService) encryptWithAES(key []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, errors.New("failed to create AES cipher: " + err.Error())
	}

	// 使用GCM模式（带认证的加密）
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, errors.New("failed to create GCM: " + err.Error())
	}

	// 生成随机IV
	iv := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, errors.New("failed to generate IV: " + err.Error())
	}

	// 加密数据
	ciphertext := aesGCM.Seal(nil, iv, plaintext, nil)

	return ciphertext, iv, nil
}

// encryptWithRSA 使用RSA-OAEP加密数据
func (s *CryptoService) encryptWithRSA(publicKeyPEM string, plaintext []byte) ([]byte, error) {
	// 解析PEM格式的公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing public key")
	}

	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse public key: " + err.Error())
	}

	rsaPubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	// 使用RSA-OAEP加密
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, plaintext, nil)
	if err != nil {
		return nil, errors.New("failed to encrypt with RSA: " + err.Error())
	}

	return ciphertext, nil
}

// DecryptMessage 解密消息（用于测试或服务端验证，实际解密在客户端）
func (s *CryptoService) DecryptMessage(privateKeyPEM string, encrypted *EncryptedMessage) (*MessageContent, error) {
	// 1. Base64解码
	encryptedAESKey, err := base64.StdEncoding.DecodeString(encrypted.EncryptedAESKey)
	if err != nil {
		return nil, errors.New("failed to decode encrypted AES key: " + err.Error())
	}

	encryptedContent, err := base64.StdEncoding.DecodeString(encrypted.EncryptedContent)
	if err != nil {
		return nil, errors.New("failed to decode encrypted content: " + err.Error())
	}

	iv, err := base64.StdEncoding.DecodeString(encrypted.IV)
	if err != nil {
		return nil, errors.New("failed to decode IV: " + err.Error())
	}

	// 2. 使用RSA私钥解密AES密钥
	aesKey, err := s.decryptWithRSA(privateKeyPEM, encryptedAESKey)
	if err != nil {
		return nil, err
	}

	// 3. 使用AES密钥解密消息内容
	plaintext, err := s.decryptWithAES(aesKey, encryptedContent, iv)
	if err != nil {
		return nil, err
	}

	// 4. 反序列化JSON
	var message MessageContent
	if err := json.Unmarshal(plaintext, &message); err != nil {
		return nil, errors.New("failed to unmarshal message: " + err.Error())
	}

	return &message, nil
}

// decryptWithRSA 使用RSA-OAEP解密数据
func (s *CryptoService) decryptWithRSA(privateKeyPEM string, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse private key: " + err.Error())
	}

	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivKey, ciphertext, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt with RSA: " + err.Error())
	}

	return plaintext, nil
}

// decryptWithAES 使用AES-GCM解密数据
func (s *CryptoService) decryptWithAES(key []byte, ciphertext []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create AES cipher: " + err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM: " + err.Error())
	}

	plaintext, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt with AES: " + err.Error())
	}

	return plaintext, nil
}
