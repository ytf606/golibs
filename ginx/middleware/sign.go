package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SignItem App签名
type SignItem struct {
	Name    string   `json:"name"`
	Source  string   `json:"source"`
	AppID   string   `json:"appId"`
	AppKey  string   `json:"appKey"`
	Expire  int64    `json:"expire"`
	Routers []string `json:"routers"`
}

// CheckSign 签名校验
func CheckSignMiddlerware(opts map[string]SignItem) gin.HandlerFunc {
	return func(c *gin.Context) {
		if opts == nil {
			c.String(http.StatusForbidden, "403 Forbidden")
			c.Abort()
			return
		}
		appSource := c.GetHeader("AUTH-SOURCE")
		if appSource == "" {
			c.String(http.StatusForbidden, "403 Forbidden")
			c.Abort()
			return
		}
		appTime := c.GetHeader("AUTH-TIME")
		if appTime == "" {
			c.String(http.StatusForbidden, "403 Forbidden")
			c.Abort()
			return
		}
		appSign := c.GetHeader("AUTH-SIGN")
		if appSign == "" {
			c.String(http.StatusForbidden, "403 Forbidden")
			c.Abort()
			return
		}
		signItem, ok := opts[appSource]
		if !ok {
			c.String(http.StatusForbidden, "403 Forbidden")
			c.Abort()
			return
		}
		iTime, err := strconv.ParseInt(appTime, 10, 64)
		if err != nil {
			c.String(http.StatusForbidden, "403 Forbidden")
			c.Abort()
			return
		}
		if iTime+signItem.Expire < time.Now().Unix() {
			c.String(http.StatusForbidden, "403 Forbidden")
			c.Abort()
			return
		}
		signBytes := md5.Sum([]byte(signItem.AppID + signItem.AppKey + "&" + appTime))
		signTxt := hex.EncodeToString(signBytes[:])
		if signTxt != appSign {
			c.String(http.StatusForbidden, "403 Forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
