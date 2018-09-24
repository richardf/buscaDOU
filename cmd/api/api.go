package main

import (
	"../../internal/db"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

type Server struct {
	db *dynamodb.DynamoDB
}

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
	srv := Server{db: db.InitDB(os.Getenv("DATABASE_URL"))}
	db.CreateStructure(srv.db)
	router := httprouter.New()
	router.POST("/keyword", srv.AddKeyword)
	router.POST("/email/:id/confirm", srv.ConfirmEmail)

	log.Fatal(http.ListenAndServe(":8080", middleware(router)))
}

func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) AddKeyword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := AddKeywordRequest{}
	json.NewDecoder(r.Body).Decode(&req)

	if ok, err := govalidator.ValidateStruct(req); ok {
		json.NewEncoder(w).Encode(&Response{Errors: err})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Response{Errors: err})
	}
}

func (s *Server) ConfirmEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ps.ByName("id")
	json.NewEncoder(w).Encode(&Response{})
}
