package v1

import (
	"fmt"
	"github.com/aenjoy/BuilderX-go/app/service"
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/hashtool"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type BaseApi struct {
}

func (f *BaseApi) GetToken(context *gin.Context) {
	cookie, err := context.Cookie("client_id")
	if err != nil {
		cookie = ioTools.Capthca2()
		context.SetCookie("client_id", cookie, 3600*24, "/", context.Request.Host, false, true)
		service.CookieInfo.Store(cookie, service.Cookie{Token: cookie, AddedTime: time.Now(), ExpiresTime: time.Now().Add(time.Hour * 24)})
		context.String(200, cookie)
	}
	context.String(http.StatusBadRequest, "cookie has been set.")
}
func (f *BaseApi) ResetCookie(context *gin.Context) {
	cookie, err := context.Cookie("client_id")
	if err != nil {
		f.GetToken(context)
		return
	}
	service.CookieInfo.Delete(cookie)
	f.GetToken(context)
}
func (f *BaseApi) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
	}
	files := form.File["files"]
	paths := form.Value["path"]
	task := form.Value["task"]
	if len(paths) == 0 || !strings.Contains(paths[0], "/") { // 末尾需要 /
		c.String(http.StatusBadRequest, "path is required")
		return
	}
	if len(task) == 0 {
		c.String(http.StatusBadRequest, "task is required")
		return
	}
	dir := path.Dir(paths[0])
	if _, err = os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("create dir err %s", err.Error()))
			return
		}
	}
	success := 0
	var hashes []string
	for _, file := range files {
		var fInfo global.FileInfo
		fInfo.Name = file.Filename
		fInfo.FullPath = path.Join("./data/upload/", paths[0], fInfo.Name)
		hash := hashtool.MD5(fInfo.FullPath + time.Now().Format("20060102150405"))
		fInfo.StorePath = path.Join("./data/upload/", paths[0], hash)
		fInfo.Task = task[0]
		if err = c.SaveUploadedFile(file, fInfo.StorePath); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err %s", err.Error()))
			return
		}
		logrus.Infoln("upload file ", fInfo.Name, " to ", fInfo.StorePath)
		logrus.Infoln("upload file ", fInfo.Name, " ok")
		success++
		global.File.Store(hash, fInfo)
		hashes = append(hashes, hash)
	}
	global.Task.Store(task[0], hashes)
	c.String(200, fmt.Sprintf("upload ok %d files", success))
}
