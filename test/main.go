package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

func GetAESEncrypted(plaintext string, secretKey string, secretIV string) (string, error) {
	var plainTextBlock []byte
	length := len(plaintext)

	extendBlock := 16 - (length % 16)
	plainTextBlock = make([]byte, length+extendBlock)
	copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))

	copy(plainTextBlock, plaintext)

	if len(plainTextBlock)%aes.BlockSize != 0 {
		fmt.Println(len(plaintext), " not equal to", aes.BlockSize)
	}
	block, err := aes.NewCipher([]byte(secretKey))

	fmt.Println(block.BlockSize(), len(secretIV))

	if err != nil {
		errorMsg := fmt.Sprintf("AES Encryption Error, val: %v, err: %v", plaintext, err.Error())
		return "", errors.New(errorMsg)
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(secretIV))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)
	return str, nil
}

func main() {
	planText := []byte("2024-2025")
	key := "kCSy7AY8yLJjqqjoiLeImKJK6Y5o3kdZ"
	secret := "mJdji*i#phiuj%^$"
	result, err := GetAESEncrypted(string(planText), key, secret)
	fmt.Println(err, result)

}
