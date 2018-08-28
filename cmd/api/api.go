package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Response struct {
}

func main() {
	router := httprouter.New()
	router.POST("/keyword", AddKeyword)
	router.POST("/email/:id/confirm", ConfirmEmail)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func AddKeyword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Welcome!")
}

func ConfirmEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s\n", ps.ByName("id"))
}
