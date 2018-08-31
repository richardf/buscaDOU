package main

import (
	"encoding/json"
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

	log.Fatal(http.ListenAndServe(":8080", middleware(router)))
}

func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

func AddKeyword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := AddKeywordRequest{}
	json.NewDecoder(r.Body).Decode(&req)

	if ok, err := govalidator.ValidateStruct(req); ok {
		json.NewEncoder(w).Encode(&Response{Errors: err})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Response{Errors: err})
	}
}

func ConfirmEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ps.ByName("id")
	json.NewEncoder(w).Encode(&Response{})
}
