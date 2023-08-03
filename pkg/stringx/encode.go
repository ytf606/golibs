package stringx

import (
	"bytes"
	"encoding/base64"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/speps/go-hashids/v2"
	"github.com/xxtea/xxtea-go/xxtea"
)

var _json = jsoniter.ConfigCompatibleWithStandardLibrary

func Decoder(data []byte, v interface{}) error {
	d := _json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	return d.Decode(&v)
}

func XxteaEncode(str string, mashIdsSalt string) string {
	encryptData := xxtea.Encrypt([]byte(str), []byte(mashIdsSalt))
	return base64.StdEncoding.EncodeToString([]byte(encryptData))
}

func XxteaDecode(str string, mashIdsSalt string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	decryptData := xxtea.Decrypt(decodeBytes, []byte(mashIdsSalt))
	return string(decryptData)
}

func Unmarshal(data []byte, v interface{}) error {
	extra.RegisterFuzzyDecoders()
	return _json.Unmarshal(data, v)
}

func UnmarshalFromString(data string, v interface{}) error {
	extra.RegisterFuzzyDecoders()
	return _json.UnmarshalFromString(data, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return _json.Marshal(v)
}

func MarshalToString(v interface{}) (string, error) {
	return _json.MarshalToString(v)
}

func MarshalAndDecode(v interface{}) (string, error) {
	encodeBytes, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encodeBytes), nil
}

func DecodeAndUnmarshal(data string, v interface{}) error {
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return Unmarshal(decodeBytes, v)
}

func MarshalAndUrlDecode(v interface{}) (string, error) {
	encodeBytes, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encodeBytes), nil
}

func UrlDecodeAndUnmarshal(data string, v interface{}) error {
	decodeBytes, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return Unmarshal(decodeBytes, v)
}

func MarshalAndRawStdDecode(v interface{}) (string, error) {
	encodeBytes, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(encodeBytes), nil
}

func RawStdDecodeAndUnmarshal(data string, v interface{}) error {
	decodeBytes, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return Unmarshal(decodeBytes, v)
}

func MarshalAndRawUrlDecode(v interface{}) (string, error) {
	encodeBytes, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encodeBytes), nil
}

func RawUrlDecodeAndUnmarshal(data string, v interface{}) error {
	decodeBytes, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return Unmarshal(decodeBytes, v)
}

func HashId(id string, mashIdsSalt string) (res string) {
	if id == "" {
		return res
	}
	hd := hashids.NewData()
	hd.Salt = mashIdsSalt
	h, _ := hashids.NewWithData(hd)

	decodeMobile, _ := h.DecodeWithError(id)
	if len(decodeMobile) == 1 {
		res = strconv.Itoa(decodeMobile[0])
	}
	return res
}
