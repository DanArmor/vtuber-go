package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/DanArmor/vtuber-go/pkg/controllers/types"
	"github.com/gin-gonic/gin"
)

func sign(payload string, key string) string {
	skHmac := hmac.New(sha256.New, []byte("WebAppData"))
	skHmac.Write([]byte(key))

	impHmac := hmac.New(sha256.New, skHmac.Sum(nil))
	impHmac.Write([]byte(payload))

	return hex.EncodeToString(impHmac.Sum(nil))
}

func checkIntegrityInitData(values url.Values, token string, expirationHours int) error {
	hash := ""
	var authDate time.Time
	pairs := make([]string, 0, len(values))
	for k, v := range values {
		if k == "hash" {
			hash = v[0]
			continue
		}
		if k == "auth_date" {
			if i, err := strconv.Atoi(v[0]); err != nil {
				authDate = time.Unix(int64(i), 0)
			}
		}
		pairs = append(pairs, k+"="+v[0])
	}
	if hash == "" {
		return errors.New("no hash")
	}
	if authDate.IsZero() {
		return errors.New("authDate is zero")
	}
	if authDate.Add(time.Duration(expirationHours)).Before(time.Now()) {
		return errors.New("authDate expired")
	}
	sort.Strings(pairs)

	if sign(strings.Join(pairs, "\n"), token) != hash {
		return errors.New("wrong hash")
	}

	return nil
}

func (s *Service) validateTgInitData(c *gin.Context) {
	var input types.InitData
	if err := c.BindQuery(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Can't bind query string"))
		return
	}
	if input.User.Id == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeNoTgId, "No Tg Id"))
		return
	}
	values := c.Request.URL.Query()
	if err := checkIntegrityInitData(values, s.TgBotToken, s.ExpirationHours); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeCantValidateInitData, "Can't validate init data"))
		return
	}

	c.Set("initData", input)
	c.Next()
}
