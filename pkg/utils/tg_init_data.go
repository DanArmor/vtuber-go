package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func sign(payload string, key string) string {
	skHmac := hmac.New(sha256.New, []byte("WebAppData"))
	skHmac.Write([]byte(key))

	impHmac := hmac.New(sha256.New, skHmac.Sum(nil))
	impHmac.Write([]byte(payload))

	return hex.EncodeToString(impHmac.Sum(nil))
}

func CheckIntegrityInitData(values url.Values, token string, expirationHours int) error {
	hash := ""
	var authDate time.Time
	pairs := make([]string, 0, len(values))
	for k, v := range values {
		if k == "hash" {
			hash = v[0]
			continue
		}
		if k == "auth_date" {
			if i, err := strconv.Atoi(v[0]); err == nil {
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
	if authDate.Add(time.Hour * time.Duration(expirationHours)).Before(time.Now()) {
		return errors.New("authDate expired")
	}
	sort.Strings(pairs)

	if sign(strings.Join(pairs, "\n"), token) != hash {
		return errors.New("wrong hash")
	}

	return nil
}
