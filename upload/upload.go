package upload

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

/////////////////*FILL THIS*//////////////////
const rootdir = ""       //root directory////
/////////////////////////////////////////////

var DBc string

///////////////////////////////////////////////////////////////////////////////////////////////
func Dashboard(w http.ResponseWriter, r *http.Request) {
	DB, err := sql.Open("mysql", DBc)
	temp, err := template.ParseFiles("templates/index.html")
	CheckErr(err)
	cookie, err := r.Cookie("ESSID")
	if err != nil || len(cookie.Value) > 20 {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if QueryCookie(cookie.Value, DB) == false {
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
		dir, err := ioutil.TempFile(rootdir, "*."+f.Filename)
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
	DB, err := sql.Open("mysql", DBc)
	CheckErr(err)
	cookie, err := r.Cookie("ESSID")
	if err != nil || len(cookie.Value) > 20 {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if QueryCookie(cookie.Value, DB) == false {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	reqfile := filepath.Join(rootdir, filepath.FromSlash(path.Clean("/"+strings.Trim(r.URL.Path[8:], "/")))) // path traversal prevention
	if filepath.Ext(reqfile) == ".html" {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeFile(w, r, reqfile)
}

///////////////////////////////////////////////////////////////////////////////////////////////

func QueryCookie(cookie string, DB *sql.DB) (stat bool) {
	stat = false
	q, err := DB.Query(`SELECT * FROM cloud WHERE session=?`, cookie)
	CheckErr(err)
	if q.Next() == true {
		stat = true
	}
	return stat
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
