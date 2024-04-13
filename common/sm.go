package common

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/emmansun/gmsm/padding"
	"github.com/emmansun/gmsm/sm3"
	"github.com/emmansun/gmsm/sm4"
)

func Sm4_d(_key string, _ciphertext string) (string, error) {
	strKey := GetSm3KeyToSm4(_key)
	key, _ := hex.DecodeString(strKey)
	ciphertext, _ := hex.DecodeString(_ciphertext)

	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < sm4.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:sm4.BlockSize]
	// fmt.Println("iv:", iv)
	ciphertext = ciphertext[sm4.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// Unpad plaintext
	pkcs7 := padding.NewPKCS7Padding(sm4.BlockSize)
	ciphertext, err = pkcs7.Unpad(ciphertext)
	if err != nil {
		return "", err
	}

	return string(ciphertext), nil
}

func Sm4_e(_key string, _ciphertext string) (string, error) {
	key, _ := hex.DecodeString(GetSm3KeyToSm4(_key))
	plaintext := []byte(_ciphertext)

	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2.
	pkcs7 := padding.NewPKCS7Padding(sm4.BlockSize)
	paddedPlainText := pkcs7.Pad(plaintext)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, sm4.BlockSize+len(paddedPlainText))
	iv := ciphertext[:sm4.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[sm4.BlockSize:], paddedPlainText)

	// fmt.Printf("%x\n", ciphertext)

	return hex.EncodeToString(ciphertext), nil

}

// 直接使用sm3.Sum方法
func GetSm3KeyToSm4(key string) string {
	sum := sm3.Sum([]byte(key))
	a := hex.EncodeToString(sum[:])
	return a[32:]
}
