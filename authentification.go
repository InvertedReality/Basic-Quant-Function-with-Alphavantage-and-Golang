package main

import (
	"./models"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var errormessage string = "POST"
var successmessage string = "User created successfully"

// POST /signup
//Create a new user
func signupAccount(writer http.ResponseWriter, request *http.Request) {
	user := models.User{}
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Println("/signup", http.StatusBadRequest, err)
		json.NewEncoder(writer).Encode("data limit exceeded")
	}
	request.Body.Close()
	if err := json.Unmarshal(body, &user); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Println("/signup", http.StatusBadRequest, err)
		json.NewEncoder(writer).Encode("Invalid json data")

	}
	if err := user.Create(); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Println("/signup", http.StatusBadRequest, err)
		json.NewEncoder(writer).Encode("Couldn't create user")
		// fmt.Println("/signup", err)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		json.NewEncoder(writer).Encode(successmessage)
		fmt.Println("/signup", http.StatusCreated)
	}
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request) {
	login_details := make(map[string]string)
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
	if err != nil {
		fmt.Println("/authenticate", http.StatusBadRequest, err)
	}
	request.Body.Close()
	if err := json.Unmarshal(body, &login_details); err != nil {
		fmt.Println("/authenticate", http.StatusBadRequest, err)
	}
	user, err := models.UserByEmail(login_details["email"])
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode("Couldn't find user")
		fmt.Println("/authenticate", http.StatusBadRequest, err)
	}
	if user.Password == models.Encrypt(login_details["password"]) {
		session, err := user.CreateSession()
		if err != nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			fmt.Println("/authenticate", http.StatusBadRequest, err)
			json.NewEncoder(writer).Encode("Error creating session")
		} else {
			expiration := time.Now().Add(time.Hour)
			cookie := http.Cookie{
				Name:     "_cookie",
				Value:    session.Uuid,
				HttpOnly: true,
				Expires:  expiration,
				MaxAge:   3600,
			}
			http.SetCookie(writer, &cookie)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(cookie)
			fmt.Println("/authenticate", http.StatusOK)
		}

	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Println("/authenticate", http.StatusBadRequest, err)
		json.NewEncoder(writer).Encode("Error verifying password")
	}

}

// GET /logout
// Logs the user out
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Failed to get cookie")
		session := models.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
}
