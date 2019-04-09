package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func Messages(status bool, message string) (map[string]interface{}) {
	return map[string]interface{} {"status" : status, "message" : message}
}

func Respond(w http.ResponseWriter, data map[string] interface{})  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

var JwtAuthentication = func(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//List of endpoints that doesn't require auth. becuause they are the initiator, it is the data they provide that we use to restrict the access
		//so the endpoints must be available
		notAuth := []string{"/api/user/new", "/api/user/login"}

		requestPath := r.URL.Path //current request path

		for _, value := range notAuth  {
			//if the request path(i.e the path that is sending the req) is among the ones we dont wat to restrict, we just send them to the next
			if value == requestPath {
				next.ServeHTTP(w,r)
				return
			}
		}
		//get the header where the token will be stored
		tokenHeader := r.Header.Get("Authentication")

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			resp := Messages( false,"There is no token in the header")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			Respond(w, resp)
			return
		}
		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement

		if len(splitted) != 2 {
			response := Messages(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}
		//extract the token from the authorization string
		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &token{}


		//this ParseWithClaims parses and validate the token, so it takes , the token that we retrieved, the claims that it will compare it with
		//and a key function which simply returns the key to be used to verify
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("lekan"), nil
		})
		if err != nil { //Malformed token, returns with http code 403 as usual
			response := Messages(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Malformed auth token", err)
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}
		if !token.Valid { //Token is invalid, maybe not signed on this server
			response := Messages(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserID + " " + tk.UserRole)
		r = r.WithContext(ctx)
		fmt.Println("context",ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!


	})

}
