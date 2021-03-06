package middleware

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	helpMessage = `
GET /log/list 查看所有日志的列表,
GET /log/list/日志名 查看对应日志的信息
GET /log/help/ 进入当前页面
POST /log/list/-1 同时使用Secret UA,上传文件`
)

func LogHelp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(helpMessage))
}

func LogList(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/log/list/")
	switch r.Method {
	case "GET":
		if id == "" {
			w.Write([]byte(ls()))
		} else {
			content, err := cat(id)
			if err != nil {
				logrus.Debug(err)
				w.Write([]byte("some thing wrong ,maybe is wrong file name\n" + helpMessage))
			}
			w.Write([]byte(content))
		}
		logrus.Debug(id)
	case "POST":
		if !rightUsrAgent(r.UserAgent()) || id != "-1" {
			w.Write([]byte(helpMessage))
			return
		}
		logrus.Info("POST log file")
		r.ParseForm()                            //解析表单
		logFile, _, err := r.FormFile("logfile") //获取文件内容

		if err != nil {
			logrus.Debug(err)
		}
		defer logFile.Close()
		logName := ""
		files := r.MultipartForm.File //获取表单中的信息
		for k, v := range files {
			for _, vv := range v {
				logrus.Println(k + ":" + vv.Filename) //获取文件名
				logName = vv.Filename
			}
			// save
		}
		saveFile, err := os.CreateTemp("", "*")

		// saveFile, err := os.T
		if err != nil {
			logrus.Debug("Create File wrong", err)
		}
		defer func() {
			newName := ("./logFile/" + logName)
			logrus.Println(newName, "logName is ", logName)
			err := os.Rename(saveFile.Name(), newName)
			if err != nil {
				logrus.Debug(err, "newFileName", newName)
			}
		}()
		defer saveFile.Close()
		_, err = io.Copy(saveFile, logFile) //保存
		if err != nil {
			logrus.Debug("save File wrong", err)
		}

		w.Write([]byte("successfully saved"))
	default:
		fmt.Println("default")
	}
}
func ls() string {
	logrus.Info("Get Log list")
	files, err := os.ReadDir("logFile")
	if err != nil {
		logrus.Debug(err)
	}
	var fileList string
	for _, v := range files {
		fileList += (v.Name() + "\n")
	}
	return fileList
}
func rightUsrAgent(ua string) bool {
	//TODO:从环境变量获取
	return ua == "vsper"
}
func cat(id string) ([]byte, error) {
	content, err := os.ReadFile("logFile/" + id)
	if err != nil {
		logrus.Debug("cat File wrong,id is", id)
	}
	return content, err
}
