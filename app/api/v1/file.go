package v1

import (
	"github.com/aenjoy/BuilderX-go/app/builder"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
)

type FileApi struct {
	dir  map[string]string //id(task) path
	init bool
}

//curl http://localhost:18088/api/v1/file/select  -X POST -d 'id=xxx&nowDir=./'

func (f *FileApi) DirList(context *gin.Context) { //目录
	if !f.init {
		f.dir = make(map[string]string)
		f.init = true
	}
	//POST
	//{id:"", nowDir:"",} id用于记录当前访问时的路径
	//outPut {File:[
	//name:"",
	//isDir:false,
	//],
	//allowBack:true,
	//status:"selectSuccess/ok/fail",
	//
	//}
	id := context.PostForm("id")
	nowDir := context.PostForm("nowDir")
	if id == "" || nowDir == "" {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	_, ok := f.dir[id]
	if !ok {
		f.dir[id] = "./" //第一次访问
	} else {
		_, err := os.Stat(path.Join(f.dir[id], nowDir))
		if err != nil {
			//路径不存在 不允许继续
			context.JSON(http.StatusBadRequest, gin.H{"status": "now dir not exist"})
			return
		}
		f.dir[id] = path.Join(f.dir[id], nowDir)
	}
	type file struct {
		Name  string `json:"name"`
		IsDir bool   `json:"isDir"`
	}
	var Files []file
	files, _ := os.ReadDir(f.dir[id])
	for _, v := range files {
		Files = append(Files, file{v.Name(), v.IsDir()})
	}
	//todo 判断是否允许上一级目录
	context.JSON(200, gin.H{"File": Files, "allowBack": true, "status": "ok"})
}

func (f *FileApi) FileSelect(context *gin.Context) { //目录或文件选择
	//{id:"", nowDir:"",}
	if !f.init {
		f.dir = make(map[string]string)
		f.init = true
	}
	id := context.PostForm("id")
	nowDir := context.PostForm("nowDir")
	if id == "" || nowDir == "" {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	info, err := os.Stat(path.Join(f.dir[id], nowDir))
	if err != nil {
		//文件不存在(一般来说不应该有)
		context.JSON(http.StatusBadRequest, gin.H{"status": "file or dir not exist. fail"})
		return
	} else {
		_, err = os.Stat(path.Join(f.dir[id], nowDir, "go.mod"))
		mod := false
		if err == nil {
			mod = true
		}
		if !info.IsDir() || mod {
			//
			t := builder.UsingAuto(path.Join(f.dir[id], nowDir), "build from web")
			if len(t) != 0 {
				builder.Tasks[id] = t[0]
				logrus.Infoln("select success. Path:", path.Join(f.dir[id], nowDir))
				context.JSON(200, gin.H{"status": "selectSuccess"})
			} else {
				context.JSON(http.StatusBadRequest, gin.H{"status": "fail. Not Config File / Dir?"})
			}
			return
		} else {
			f.DirList(context)
			return
		}
	}

}
