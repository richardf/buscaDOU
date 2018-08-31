package main

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Response struct {
	Errors error
}

type AddKeywordRequest struct {
	Keyword string `json:"keyword" valid:"length(5|64)"`
	Email   string `json:"email" valid:"email"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func main() {
	router := httprouter.New()
	router.POST("/keyword", AddKeyword)
	router.POST("/email/:id/confirm", ConfirmEmail)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func AddKeyword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := AddKeywordRequest{}
	json.NewDecoder(r.Body).Decode(&req)
	w.Header().Set("Content-Type", "application/json")

	if ok, err := govalidator.ValidateStruct(req); ok {
		json.NewEncoder(w).Encode(&Response{Errors: err})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Response{Errors: err})
	}
}

func ConfirmEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s\n", ps.ByName("id"))
}
