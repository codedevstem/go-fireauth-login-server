package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

const port = ":13300"

// RegisterUser is the struct that is expected when a user register for the application
type RegisterUser struct {
	Email    string `json:"email" validate:"required,email,max=60"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

var validate *validator.Validate

func registerHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		validate := validator.New()

		registerUser := RegisterUser{}

		err := json.NewDecoder(request.Body).Decode(&registerUser)

		if err != nil {
			writer.WriteHeader(400)
			fmt.Fprint(writer, "Body is not valid")
		}

		err = validate.Struct(registerUser)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				writer.WriteHeader(400)
				fmt.Printf("A validation error occured: %v", err)
				fmt.Fprint(writer, err)

			}
		}

		fmt.Fprint(writer, "Request method not supported.")
	}
	writer.WriteHeader(200)
	fmt.Fprint(writer, "Request OK")

}

func loginHandler(writer http.ResponseWriter, request *http.Request) {

}

func logoutHandler(writer http.ResponseWriter, request *http.Request) {

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/logout", logoutHandler)

	http.ListenAndServe(port, mux)

}
