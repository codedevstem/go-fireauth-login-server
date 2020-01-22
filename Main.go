package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/go-playground/validator.v9"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

const port = ":13300"

// RegisterUser is the struct that is expected when a user register for the application
type RegisterUser struct {
	Email    string `json:"email" validate:"required,email,max=60"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginUser is the struct that is expected when a user preforms a login
type LoginUser struct {
	Email    string `json:"Email" validate:"required,email,max=60"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreatedUser struct {
	idToken      string
	email        string
	refreshToken string
	expiresIn    string
	localId      string
}

var validate *validator.Validate

var app *firebase.App

var firebaseAPIKey string

const firebaseRestBasePath = "https://identitytoolkit.googleapis.com/v1"
const accountPath = "/accounts"

func registerHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

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
		url := firebaseRestBasePath + accountPath + ":signUp?key=" + firebaseAPIKey
		jsonBody, _ := json.Marshal(registerUser)
		body := bytes.NewBuffer(jsonBody)
		resp, err := http.Post(url, "application/json", body)

		if err != nil {
			fmt.Fprintf(writer, "Failed to create user: %v", err)
		}

		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Firebase response not parsable.")
			fmt.Fprint(writer, "User not created")
		}
		defer resp.Body.Close()

		fmt.Printf("Created user: %s", registerUser.Username)

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(responseData)
	} else {
		fmt.Fprint(writer, "Request method not supported.")
	}
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	loginUser := LoginUser{}

	err := json.NewDecoder(request.Body).Decode(&loginUser)
	if err != nil {
		fmt.Fprint(writer, "Body is not valid")
	}

	err = validate.Struct(loginUser)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			writer.WriteHeader(400)
			fmt.Printf("A validation error occured: %v", err)
			fmt.Fprint(writer, err)
		}
	}

}

func logoutHandler(writer http.ResponseWriter, request *http.Request) {

}

func main() {
	mux := http.NewServeMux()

	content, err := ioutil.ReadFile("credentials/todos-app-fdef3-api-key")
	if err != nil {
		log.Fatal(err)
	}

	firebaseAPIKey = string(content)

	opt := option.WithCredentialsFile("credentials/todos-app-fdef3-service-account-key.json")
	app, err = firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		fmt.Printf("error initializing app: %v", err)
		os.Exit(-1)
	}

	validate = validator.New()

	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/logout", logoutHandler)

	http.ListenAndServe(port, mux)

}
