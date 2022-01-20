package main

import (
	"evpeople/ChatRoom/middleware"
	"evpeople/ChatRoom/ws"
	"flag"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

var addr = flag.String("addr", ":8081", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.Trace("begin working")
	flag.Parse()
	hub := ws.NewHub()
	//middleware.Meth("Get")的执行结果是一个函数闭包，这个函数存储了"GET"的信息
	http.HandleFunc("/login", middleware.Login)
	http.HandleFunc("/ss", middleware.Chain(middleware.Hello, middleware.Method("GET"), middleware.CheckCookie()))
	// http.HandleFunc("/login")
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	}, middleware.CheckCookie()))
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	ws.ServeWs(hub, w, r)
	// })
	go hub.Run()
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
