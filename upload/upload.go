package upload

import (
	"github.com/go-redis/redis"
	"html/template"
	"io/ioutil"
	"log"
	"logger"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

/////////////////*FILL THIS*//////////////////
const rootdir = ""       //root directory////
/////////////////////////////////////////////

var DBc string
var rediscnt redis.Options = redis.Options{Addr: os.Getenv("REDISADDR") + ":6379", Password: os.Getenv("REDISPWD"), DB: 0}

///////////////////////////////////////////////////////////////////////////////////////////////
func Dashboard(w http.ResponseWriter, r *http.Request) {
	logger.LogInsert(r.Method, r.UserAgent(), r.URL.Path, r)
	temp, err := template.ParseFiles("templates/index.html")
	CheckErr(err)
	cookie, err := r.Cookie("ESSID")
	if err != nil || len(cookie.Value) > 20 {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if QueryCookie(cookie.Value) == false {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if r.Method == "GET" {
		temp.Execute(w, nil)
		return
	}
	if r.Method == "POST" {
		r.ParseMultipartForm(100000000)
		file, f, err := r.FormFile("file")
		if err != nil {
			temp.Execute(w, "Invalid File")
			return
		}
		defer file.Close()
		if f.Size > 100000000 {
			temp.Execute(w, "Error. file is too large")
			return
		}
		dir, err := ioutil.TempFile(rootdir, "*-"+f.Filename)
		CheckErr(err)
		defer dir.Close()
		fb, err := ioutil.ReadAll(file)
		CheckErr(err)
		_, err = dir.Write(fb)
		CheckErr(err)
		temp.Execute(w, f.Filename+" uploaded successfully!")
	}
}

func Uploads(w http.ResponseWriter, r *http.Request) {
	logger.LogInsert(r.Method, r.UserAgent(), r.URL.Path, r)
	cookie, err := r.Cookie("ESSID")
	if err != nil || len(cookie.Value) > 20 {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if QueryCookie(cookie.Value) == false {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	reqfile := filepath.Join(rootdir, filepath.FromSlash(path.Clean("/"+strings.Trim(r.URL.Path[8:], "/")))) // path traversal prevention
	if CheckFile(filepath.Ext(reqfile)) == false {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeFile(w, r, reqfile)
}

///////////////////////////////////////////////////////////////////////////////////////////////

func QueryCookie(cookie string) (stat bool) {
	rdb := redis.NewClient(&rediscnt)
	_, err := rdb.Get(cookie).Result()
	if err == redis.Nil {
		stat = false
	} else {
		stat = true
	}
	return stat
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func CheckFile(ext string) bool {
	legalfiles := [9]string{".png", ".jpeg", ".jpg", ".gif", ".pdf", ".mp4", ".mov", ".mp3", ".txt"}
	arr := reflect.ValueOf(legalfiles)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == ext {
			return true
		}
	}
	return false
}
