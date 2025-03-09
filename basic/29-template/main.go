package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type User struct {
	Name    string
	Age     int
	Hobbies []string
}

func main() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {

		// 准备数据
		user := User{
			Name:    "张三",
			Age:     18,
			Hobbies: []string{"篮球", "足球", "羽毛球"},
		}

		// 解析模版
		tmpl, err := template.ParseFiles("index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 渲染模版
		err = tmpl.Execute(w, user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	// http.ListenAndServe(":8080", nil)

	ticker := time.Tick(time.Second)

	for i := range ticker {
		fmt.Println(i)
	}
}
