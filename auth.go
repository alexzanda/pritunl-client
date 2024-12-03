package pritunl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 生成认证的请求头, requestPath样例/server，不能包含协议和主机信息部分，且不应包含查询参数部分
func generateAuthHeader(requestPath, requestMethod, apiToken, apiSecret string) (map[string]string, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	authTimestamp := strconv.Itoa(int(time.Now().Unix()))

	// authNonce生成，最长只能为32
	authNonce := strings.ReplaceAll(uuid.New().String(), "-", "")

	// signature生成
	authString := strings.Join([]string{apiToken, authTimestamp, authNonce, strings.ToUpper(requestMethod), requestPath}, "&")
	signature, err := generateHMACBase64(apiSecret, authString)
	if err != nil {
		return nil, err
	}

	// 认证信息组装
	header["Auth-Token"] = apiToken
	header["Auth-Timestamp"] = authTimestamp
	header["Auth-Nonce"] = authNonce
	header["Auth-Signature"] = signature
	return header, nil
}

// 生成signature
func generateHMACBase64(secret, message string) (string, error) {
	// 创建一个新的 HMAC-SHA256 哈希器
	h := hmac.New(sha256.New, []byte(secret))

	if _, err := h.Write([]byte(message)); err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return encoded, nil
}
