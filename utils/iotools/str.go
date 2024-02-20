package iotools

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 生成6位随机验证码（数字）(生成6位ID)
func Captcha1() (int, error) {
	a, b := strconv.Atoi(fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000)))
	logrus.Info("Captcha1Api:", a, b)
	return a, b
}

// 生成16位随机ID（字母）
func Capthca2() string {
	n := 16
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	logrus.Info("Captcha2Api:", sb.String())
	return sb.String()
}

func IsStrAInStrB(a string, b string) bool {
	return strings.Contains(b, a)
}
