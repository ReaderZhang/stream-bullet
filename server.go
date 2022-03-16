package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)
import "stream-bullet/server"

func main() {

	router := mux.NewRouter()
	//开启协程启动connection服务管理
	go server.H.Run()
	//创建ws服务
	router.HandleFunc("/ws", server.Myws)
	//启动http服务
	if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
		fmt.Println("err:", err)
	}
}
