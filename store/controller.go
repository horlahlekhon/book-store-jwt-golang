package store

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)


func Message(status bool, message string) map[string]interface{}  {
	return map[string]interface{}{"status": status, "message": message}
}

var ServeBooks = func(w http.ResponseWriter, r *http.Request) {

	s := r.Context().Value("user")

	fmt.Println("context variable in ServeBooks", s)

	data := GetBooks()
	body := map[string]interface{}{
		"data" : data,
		//"Message" : "Thanks for banking with us",
	}

	w.Header().Add("Content-Type", "application/json")
	//set the response status header
	w.WriteHeader(http.StatusOK)

	//convert the body we want to send to byte
	 by, error := json.Marshal(&body)
	if error != nil {
		fmt.Println(error)
	}

	 //write the byte to the response writer
	 w.Write(by)

}

// TODO this is not working properly, there is something wrong with the json decoding the passed in object
var AddBook = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	ctx := r.Context().Value("user")
	fmt.Println("user ::",ctx)
	book := &Book{}

	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(error.Error())
		return
	}
	er := r.Body.Close()
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(er)
		return
	}

	err :=  json.Unmarshal(body, book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	bok := book.CreateBook()
	resp := Message(true, "success")
	resp["data"] = bok


	data , err := json.Marshal(resp)

	w.Write(data)
}

var GetBookById = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	//id := r.Context().Value("id") . (int)
	fmt.Println("book id", params["id"])

	id, error := strconv.Atoi(params["id"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(error.Error())
		return

	}

	encodedBook := GetBook(id)

	w.Header().Add("Content-Type", "application/json")

	 body,er := json.Marshal(encodedBook)
	 if er != nil {
	 	fmt.Println(er.Error())
	 	w.WriteHeader(http.StatusInternalServerError)
		 return
	 }

	 w.Write(body)
}

var PatchBook = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	body , error := ioutil.ReadAll(r.Body)

	if error !=nil {
		fmt.Println(error)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := r.Body.Close()
	if err != nil {
		fmt.Println(error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var book Book
	//convert the body from byte to the interface specified
	erro := json.Unmarshal(body, &book)
	if erro != nil{
		fmt.Println(erro.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	retBook := UpdateBook(book)

	updateResp , erru := json.Marshal(retBook)

	if erru!= nil {
		fmt.Println(err)
		panic(err)
		}


	w.Write(updateResp)

}

var DeleteBook = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, error := strconv.Atoi(params["id"])
	if error != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(error.Error())
		return
	}
	data := DeleteBookById(id)

	w.Header().Add("Content-Type", "application/json")
	body, er := json.Marshal(data)
	if er != nil {
		fmt.Println(er.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

// restricted

type Credentials struct {
	Name string `json:"name"`
	Password string `json:"password"`
}

var Logon =  func(w http.ResponseWriter, r *http.Request) {

				w.Header().Add("Content-Type", "application/json")
				ctxx := r.Context().Value("user")
				fmt.Println("user id",ctxx)

				body , error := ioutil.ReadAll(r.Body)
				if error != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				err := r.Body.Close()
				if err != nil {
					fmt.Println(error)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				var credentials Credentials
				erro := json.Unmarshal(body, &credentials)
				if erro != nil{
					fmt.Println(erro.Error())
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				boolean, data, id := Login(credentials.Name, credentials.Password)

				if !boolean {
					fmt.Println( "", data)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error" : "No User with that credentials"}`))
					return
				}


				w.Header().Add("Authentication", "Bearer " + data)



				bo , er := json.Marshal("logged in successfully " + id + " " + "token: " + data  )
				if er!=nil{
					fmt.Println(erro.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}


				w.Write(bo)


}
var Register = func(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	body , error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(error.Error())
		return
	}

	var account Account
	erro := json.Unmarshal(body, &account)
	if erro != nil{
		fmt.Println(erro.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("account in controller registier", account)
	boolean , data := RegisterAccount(account)
	if !boolean {
		fmt.Println(data)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, er := json.Marshal(data)

	if er!=nil{
		fmt.Println(erro.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

var GetUsers = func(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	user := r.Context().Value("user")
	if user == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		by, error := json.Marshal("Invalid authentication token")
		if error != nil {
			fmt.Println(error)
		}
		w.Write(by)
	}else {

		data := user.(string)
		splitted := strings.Split(data," ")

		if len(splitted) != 2 {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			by, error := json.Marshal("Invalid authentication token")
			if error != nil {
				fmt.Println(error)
			}
			w.Write(by)
		}else {
			if splitted[1] == "admin" {
				data := GetAllAccounts()
				w.Header().Add("Content-Type", "application/json")
				//set the response status header
				w.WriteHeader(http.StatusOK)

				//convert the body we want to send to byte
				by, error := json.Marshal(&data)
				if error != nil {
					fmt.Println(error)
				}

				//write the byte to the response writer
				w.Write(by)

			}else {

				data, err := GetAccount(splitted[0])

				if err != "" {
					fmt.Println("wait! who gave you this token, %v", err)
				}
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				by, error := json.Marshal(&data)
				if error != nil {
					fmt.Println(error)
				}
				w.Write(by)
			}
		}


	}


}