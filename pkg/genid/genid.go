package genid

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/spf13/cast"
)

func CalcMd5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func RandStringUUid() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	timeNow := time.Now().Unix()
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return CalcMd5Hash(string(b) + cast.ToString(timeNow))[0:25]
}

func GenId() (string, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", err
	}
	return node.Generate().String(), nil
}
