package main

import (
	"./models"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	// "strings"
	"github.com/dgrijalva/jwt-go"
)

var successmessage string = "User created successfully"

// GET  /userlogout
// POST /user/auth/

// Authenticate the user given the email and password and logout
func authenticate(writer http.ResponseWriter, request *http.Request) {
	switch{
	case request.Method=="POST" && request.URL.Path=="/user/auth/" :
		{
			login_details := make(map[string]string)
			body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
			if err != nil {
				fmt.Println("/user/auth/", http.StatusBadRequest, err)
			}
			request.Body.Close()
			if err := json.Unmarshal(body, &login_details); err != nil {
				fmt.Println("/user/auth/", http.StatusBadRequest, err)
			}
			user, err := models.UserByEmail(login_details["email"])
			if err != nil {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(writer).Encode("Couldn't find user")
				fmt.Println("/user/auth/", http.StatusBadRequest, err)
			}
			if user.Password == models.Encrypt(login_details["password"]) {
				if err != nil {
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusBadRequest)
					fmt.Println("/user/auth/", http.StatusBadRequest, err)
					json.NewEncoder(writer).Encode("Error creating session")
				} else {
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusOK)
					fmt.Println("/user/auth/", http.StatusOK)
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			        "user": user,
			        "exp":time.Now().Add(time.Hour * 2),

			    })
			    tokenString, error := token.SignedString([]byte("secret"))
			    if error != nil {
			        fmt.Println(error)
			    }
			    json.NewEncoder(writer).Encode(JwtToken{Token: tokenString})
				}

			} else {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusBadRequest)
				fmt.Println("/user/auth/", http.StatusBadRequest, err)
				json.NewEncoder(writer).Encode("Error verifying password")
			}

		}
	case request.Method=="GET" && request.URL.Path=="/user/logout/":
		{
			cookie, err := request.Cookie("_cookie")
			if err != http.ErrNoCookie {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(400)
				json.NewEncoder(writer).Encode("Failed to get cookie")
				session := models.Session{Uuid: cookie.Value}
				session.DeleteByUUID()
			}
		}
	}
}
// POST /user/signup/
//GET /user/list/
//POST /user/delete/
func UserExec(writer http.ResponseWriter, request *http.Request){
	switch{
	case request.Method=="GET" && request.URL.Path=="/user/list/" :
		 {
			users,err:=models.Users()
			if err != nil {
				{
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Println(http.StatusInternalServerError)
					fmt.Println(err)
				}
			} else {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				encoder := json.NewEncoder(writer)
				encoder.SetIndent(empty, tab)
				encoder.Encode(users)
				fmt.Println("/user/list/", http.StatusOK)
			}
		}
	case request.Method=="POST" && request.URL.Path=="/user/delete/" :
		{
			user_email := make(map[string]string)

			body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))

			if err != nil {
				fmt.Println("/user/delete/", http.StatusBadRequest, err)
			}
			request.Body.Close()
			if err := json.Unmarshal(body, &user_email); err != nil {
				fmt.Println("/user/delete/", http.StatusBadRequest, err)
			}else{
				user, err := models.UserByEmail(user_email["email"])
				if err != nil {
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(writer).Encode("Couldn't find user")
					fmt.Println("/user/delete/", http.StatusBadRequest, err)
				}else{
					user.Delete()
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusOK)
					fmt.Println("/user/delete/", http.StatusOK)
				}
			}
		}
	case request.Method=="POST" && request.URL.Path=="/user/signup/":
		{
			user := models.User{}
			body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
			if err != nil {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusBadRequest)
				fmt.Println("/user/signup/", http.StatusBadRequest, err)
				json.NewEncoder(writer).Encode("data limit exceeded")
			}
			request.Body.Close()
			if err := json.Unmarshal(body, &user); err != nil {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusBadRequest)
				fmt.Println("/user/signup/", http.StatusBadRequest, err)
				json.NewEncoder(writer).Encode("Invalid json data")

			}
			if err := user.Create(); err != nil {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusBadRequest)
				fmt.Println("/user/signup/", http.StatusBadRequest, err)
				json.NewEncoder(writer).Encode("Couldn't create user")
				
			} else {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusCreated)
				json.NewEncoder(writer).Encode(successmessage)
				fmt.Println("/user/signup/", http.StatusCreated)
			}
		}
	}
}
