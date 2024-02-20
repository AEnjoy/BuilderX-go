package hashtool

import (
	"crypto/md5"
	"fmt"
	"github.com/sirupsen/logrus"
)

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	logrus.Debug("hash str is ", str, ",and it MD5", md5str)
	return md5str
}
