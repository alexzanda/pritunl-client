package pritunl

import (
	"fmt"
	"math/rand"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitBytes  = "0123456789"
	allBytes    = letterBytes + digitBytes
)

// generateVpnServerName 生成一个指定长度的随机字符串
func generateVpnServerName() string {
	length := 5

	b := make([]byte, length)
	for i := 1; i < length; i++ {
		b[i] = allBytes[rand.Intn(len(allBytes))]
	}

	return fmt.Sprintf("server_%s", string(b))
}
