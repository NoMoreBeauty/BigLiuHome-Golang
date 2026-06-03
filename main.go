package main

import (
	"fmt"
	"log"
	"net/http"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"

	"github.com/gorilla/mux" // 引入 mux
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	r := mux.NewRouter()

	// 静态路由
	r.HandleFunc("/", service.IndexHandler).Methods("GET")
	r.HandleFunc("/api/count", service.CounterHandler).Methods("GET", "POST")
	r.HandleFunc("/api/login", service.AuthHandler).Methods("GET")
	r.HandleFunc("/api/meals", service.GetMealsHandler).Methods("GET", "POST")
	r.HandleFunc("/api/meals/calendar", service.GetMealsCalendarHandler).Methods("GET")

	// 动态路由
	r.HandleFunc("/api/meals/{id:[0-9]+}", service.GetMealDetailHandler).Methods("GET")
	r.HandleFunc("/api/meals/{id:[0-9]+}/comments", service.GetMealCommentsHandler).Methods("GET")
	r.HandleFunc("/api/meals/{id:[0-9]+}/comments", service.PostCommentHandler).Methods("POST")
	r.HandleFunc("/api/meals/{id:[0-9]+}/like", service.PostLikeHandler).Methods("POST")

	// 使用 r 作为 HTTP handler 启动服务
	log.Fatal(http.ListenAndServe(":80", r))
}
