package filecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func Encrypt(file string, password []byte) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		panic(err.Error())
	}

	srcFile, err := os.Open(file)

	if err != nil {
		panic(err.Error())
	}

	defer srcFile.Close()

	plainText, err := io.ReadAll(srcFile)

	if err != nil {
		panic(err.Error())
	}

	key := password
	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha1.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	cipherText := aesGCM.Seal(nil, nonce, plainText, nil)
	cipherText = append(cipherText, nonce...)

	dstFile, err := os.Create(file)
	if err != nil {
		panic(err.Error())
	}

	defer dstFile.Close()

	_, err = dstFile.Write(cipherText)
	if err != nil {
		panic(err.Error())
	}

}

func Decrypt(file string, password []byte) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		panic(err.Error())
	}

	srcFile, err := os.Open(file)
	if err != nil {
		panic(err.Error())
	}

	defer srcFile.Close()

	cipherText, err := io.ReadAll(srcFile)
	if err != nil {
		panic(err.Error())
	}

	key := password
	salt := cipherText[len(cipherText)-12:]
	str := hex.EncodeToString(salt)
	nonce, err := hex.DecodeString(str)
	if err != nil {
		panic(err.Error())
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha1.New)
	block, err := aes.NewCipher(dk)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plainText, err := aesGCM.Open(nil, nonce, cipherText[:len(cipherText)-12], nil)
	if err != nil {
		panic(err.Error())
	}

	dstFile, err := os.Create(file)
	if err != nil {
		panic(err.Error())
	}

	defer dstFile.Close()

	_, err = dstFile.Write(plainText)
	if err != nil {
		panic(err.Error())
	}
}
