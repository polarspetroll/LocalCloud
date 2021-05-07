package main

import (
	"fmt"
	"log"
	"login"
	"net/http"
	"os"
	"upload"
)

func main() {
	DB := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", os.Getenv("DBUSR"), os.Getenv("DBPWD"), os.Getenv("DBADDR"), os.Getenv("DBNAME"))
	login.DBc = DB
	upload.DBc = DB
	statics := http.FileServer(http.Dir("./statics"))
	//////////////////////////////////////////////////////////////////////
	http.Handle("/statics/", http.StripPrefix("/statics/", statics))
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/login", login.Login)
	http.HandleFunc("/dashboard", upload.Dashboard)
	http.HandleFunc("/changepassword", login.ChangePassword)
	http.HandleFunc("/uploads/", upload.Uploads)
	//////////////////////////////////////////////////////////////////////
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard", 302)
}
