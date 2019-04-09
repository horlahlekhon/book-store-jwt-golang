package main

import (
	"book-store/store"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)


func main() {


	router := mux.NewRouter()

	router.Use(store.JwtAuthentication)
	router.HandleFunc("/api/user/new", store.Register).Methods("POST")
	router.HandleFunc("/api/user/login", store.Logon).Methods("POST")
	router.HandleFunc("/api/book/{id}/", store.GetBookById).Methods("GET")
	router.HandleFunc("/api/book/add/", store.AddBook).Methods("POST")
	router.HandleFunc("/api/books/", store.ServeBooks).Methods("GET")
	router.HandleFunc("/api/book/update/", store.PatchBook).Methods("PATCH")
	router.HandleFunc("/api/book/remove/{id}/", store.DeleteBook).Methods("DELETE")
	router.HandleFunc("/api/users", store.GetUsers).Methods("GET")




	err := http.ListenAndServe(":9000", router)


	if err != nil {
		fmt.Println(err)
	}
}
