package login

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
)

var DBc string

func Login(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/login.html")
	CheckErr(err)
	DB, err := sql.Open("mysql", DBc)
	CheckErr(err)
	if r.Method == "GET" {
		tmp.Execute(w, nil)
		return
	} else if r.Method == "POST" {
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if len(username) > 40 {
			tmp.Execute(w, "Invalid Username")
			return
		}
		stat, sid := Query(username, password, DB)
		if stat == false && sid == "" {
			tmp.Execute(w, "incorrect username or password")
			return
		} else if stat == false && sid != "" {
			tmp.Execute(w, "Server Error")
			return
		}
		ck := http.Cookie{Name: "ESSID", Value: sid}
		http.SetCookie(w, &ck)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/passwordchange.html")
	CheckErr(err)
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
	if r.Method == "GET" {
		tmp.Execute(w, nil)
		return
	} else if r.Method == "POST" {
		username := r.PostFormValue("username")
		if len(username) > 40 {
			tmp.Execute(w, "invalid username")
			return
		}
		password := r.PostFormValue("password")
		confirm := r.PostFormValue("confirm")
		if password != confirm {
			tmp.Execute(w, "passwords does not match")
			return
		}
		stat := PasswordUpdate(username, password, DB)
		if stat != 1 {
			tmp.Execute(w, "Invalid Username")
			return
		}
		tmp.Execute(w, "Password Changed Successfully")
	} else {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
}

////////////////
func Query(username, password string, DB *sql.DB) (stat bool, sid string) {
	stat = true
	q, err := DB.Query(`SELECT username, password FROM cloud WHERE username=? AND password=MD5(?)`, username, password)
	CheckErr(err)
	if q.Next() == true {
		a := make([]byte, 10)
		rand.Read(a)
		sid = hex.EncodeToString(a)
		if Update(sid, username, DB) != 1 {
			stat = false
		}
	} else if q.Next() == false {
		sid = ""
		stat = false
	}
	return stat, sid
}

////////////////
func Update(sid, username string, DB *sql.DB) int64 {
	p, err := DB.Prepare(`UPDATE cloud SET session=? WHERE username=?`)
	CheckErr(err)
	res, err := p.Exec(sid, username)
	CheckErr(err)
	row, err := res.RowsAffected()
	CheckErr(err)
	return row
}

func QueryCookie(cookie string, DB *sql.DB) (stat bool) {
	stat = false
	q, err := DB.Query(`SELECT session FROM cloud WHERE session=?`, cookie)
	CheckErr(err)
	if q.Next() == true {
		stat = true
	}
	return stat
}

func PasswordUpdate(username, password string, DB *sql.DB) int64 {
	q, err := DB.Prepare(`UPDATE cloud SET password=MD5(?) WHERE username=? `)
	CheckErr(err)
	e, err := q.Exec(password, username)
	CheckErr(err)
	row, err := e.RowsAffected()
	CheckErr(err)
	return row
}

////////////////

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
