package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	rpcHelper "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
)

var (
	encryptionKey []byte
	UseEncryption bool
	logging       rpcHelper.LoggingClient
)

func SetLogging(client rpcHelper.LoggingClient) {
	logging = client
}

func SetupCrypt(key []byte) {
	encryptionKey = key

	UseEncryption = true // keep this as FALSE for the moment. Needs to become TRUE (and/or removed from code) in V32.
}

func EncryptRaw(data []byte) []byte {
	block, _ := aes.NewCipher(encryptionKey)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

// Converts a string to HEX after Encryption.
func Encrypt(data string) string {
	encoded := EncryptRaw([]byte(data))
	encodedHex := hex.EncodeToString(encoded)
	return encodedHex
}

func DecryptRaw(data []byte) ([]byte, error) {
	if len(data) <= 0 {
		return data, nil
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("Data too short for decryption: " + hex.EncodeToString(data))
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func Decrypt(data string) (string, error) {
	hexCoded, err := hex.DecodeString(data)
	if err != nil {
		if logging != nil {
			logging.Error(err)
		} else {
			//This should at least get it in journalctl
			println(err.Error())
		}
		return "", err
	}

	decoded, err := DecryptRaw(hexCoded)
	if err != nil {
		if logging != nil {
			logging.Error(err)
		} else {
			//This should at least get it in journalctl
			println(err.Error())
		}
		return "", err
	}

	return string(decoded), nil
}
