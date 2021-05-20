package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// トップページ
	http.Handle("/", http.FileServer(http.Dir("static")))
	// プロセスを起動
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
