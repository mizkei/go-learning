package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gotschmarcel/goserv"
)

// User is not special
type User struct {
	name string
}

func nothingHandler(w http.ResponseWriter, r *http.Request) {
	goserv.WriteString(w, "nothing")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("認証したことにしよう")
	goserv.Context(r).Set("user", User{name: "guest"})
}

func afterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("処理終わったよ")
}

func simpleHandler(str string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		goserv.WriteString(w, str)
		log.Println("response:", str)
	}
}

func main() {
	server := goserv.NewServer()

	server.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err *goserv.ContextError) {
		w.WriteHeader(err.Code)
		goserv.WriteString(w, err.String())
		log.Println("error handling")
	}

	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		goserv.WriteString(w, "It works")
	})

	// Routeでメソッドチェーンを書くことができる
	server.Route("/articles").Get(nothingHandler).Post(nothingHandler)

	// apiとかのprefixつけるのも簡単
	api := server.SubRouter("/api")
	api.Use(authHandler) // 認証などの処理はUseですべてのapiに適用

	api.Get("/bookmark", simpleHandler("GET bookmark"))
	api.Post("/bookmark", simpleHandler("POST bookmark"))
	api.Get("/bookmark/:article_id", func(w http.ResponseWriter, r *http.Request) {
		id := goserv.Context(r).Param("article_id")

		if id == "nanikore" {
			goserv.Context(r).Error(fmt.Errorf("invalid id"), 400)
			return
		}

		goserv.WriteString(w, "article id:"+id)
		log.Println("GET bookmark/:article_id")
	})
	api.Param("article_id", func(w http.ResponseWriter, r *http.Request, id string) {
		if id == "invalid" {
			goserv.WriteString(w, "invalidなやつ")
			return
		}
	})

	api.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		user := goserv.Context(r).Get("user").(User)
		goserv.WriteString(w, "user: "+user.name)
	})
	api.Get("/article", func(w http.ResponseWriter, r *http.Request) {
		var article struct{ Title string }
		article.Title = "Lobi"

		goserv.WriteJSON(w, &article)
	})

	server.Use(afterHandler) // Useは順番も重要。これは何も処理が行われなかった時に適用される

	log.Fatalln(server.Listen(":8888"))
}
