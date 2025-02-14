package main

import (
	"cmp"
	"log"
	"net/http"
	"os"

	"github.com/crhntr/dom/domtest/example/blog"
)

func main() {
	mux := http.NewServeMux()
	app := new(blog.App)
	blog.Routes(mux, app)
	log.Fatal(http.ListenAndServe(":"+cmp.Or(os.Getenv("PORT"), "8080"), mux))
}
