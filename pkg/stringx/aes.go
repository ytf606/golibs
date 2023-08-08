package stringx

import (
	"crypto/aes"
	"encoding/base64"

	"github.com/ytf606/golibs/logx"
	"github.com/ytf606/golibs/pkg/osx"
)

func AesEncryptECB(origDataStr, keyEcb string) string {
	origData := []byte(origDataStr)
	key := []byte(keyEcb)
	cipher, err := aes.NewCipher(generateKey(key))
	if err != nil {
		logx.E(osx.PF(), "mobile aes encode failed raw mobile:%s, err:%+v", origDataStr, err)
		return ""
	}

	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}

	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return base64.StdEncoding.EncodeToString(encrypted)
}

func AesDecryptECB(encryptedStr, keyEcb string) string {
	key := []byte(keyEcb)
	encryptedDecode, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		logx.E(osx.PF(), "mobile aes decode failed raw mobile:%s, err:%+v", encryptedStr, err)
		return ""
	}

	encrypted := encryptedDecode
	cipher, err := aes.NewCipher(generateKey(key))
	if err != nil {
		logx.E(osx.PF(), "mobile aes decode failed base64 mobile:%s, err:%+v", encryptedStr, err)
		return ""
	}

	decrypted := make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return string(decrypted[:trim])
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
