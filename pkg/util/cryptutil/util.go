package cryptutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// CreateHash : Hash a string with md5
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// EncryptData : Encrypt given data with a password
func EncryptData(data []byte, passwword string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(createHash(passwword)))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// DecryptData : Decrypt given data with a password
func DecryptData(data []byte, password string) ([]byte, error) {
	key := []byte(createHash(password))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// EncryptStringToBase64 : Encrypt given string with a password to an base64 string
func EncryptStringToBase64(data string, password string) (string, error) {
	encyrptData, err := EncryptData([]byte(data), password)
	if err != nil {
		return "", err
	}
	base64Data := base64.StdEncoding.EncodeToString(encyrptData)
	return base64Data, nil
}

// DecryptBase64ToString : Decrypt given base64 string with a password to an string
func DecryptBase64ToString(data string, password string) (string, error) {
	base64Decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	decryptedData, err := DecryptData([]byte(base64Decoded), password)
	if err != nil {
		return "", err
	}
	return string(decryptedData), nil
}
