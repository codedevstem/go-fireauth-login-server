package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"

	"github.com/codedevstem/go-fireauth-login-server/register"
)

const port = ":13300"

// LoginUser is the struct that is expected when a user preforms a login
type LoginUser struct {
	Email    string `json:"Email" validate:"required,email,max=60"`
	Password string `json:"password" validate:"required,min=8"`
}


func logoutHandler(writer http.ResponseWriter, request *http.Request) {

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", register.Handler)
	mux.HandleFunc("/logout", logoutHandler)

	http.ListenAndServe(port, mux)

}
