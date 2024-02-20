package debugTools

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var StartTime = time.Now()

func init() {
	t := StartTime.Format("2006-01-02")
	logf, _ := os.OpenFile("./logs/"+t+"iomServer.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	mw := io.MultiWriter(os.Stdout, logf)
	logrus.SetOutput(mw)
}
