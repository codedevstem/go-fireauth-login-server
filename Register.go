package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

// UserToRegister is the struct that is expected when a user register for the application
type UserToRegister struct {
	Email    string `json:"email" validate:"required,email,max=60"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

// UserThatsRegistered is the struct that is expected from the register API in firebase
type UserThatsRegistered struct {
	idToken			string
	email			string
	refreshToken	string
	expiresIn		string
	localId			string	
}

const firebaseRestBasePath = "https://identitytoolkit.googleapis.com/v1"
const accountPath = "/accounts"

var validate *validator.Validate

var firebaseAPIKey string

func init() {
	validate = validator.New()
	content, err := ioutil.ReadFile("credentials/todos-app-fdef3-api-key")
	if err != nil {
		log.Fatal(err)
	}
	firebaseAPIKey = string(content)
}

func Handler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		userToRegister := UserToRegister{}

		err := json.NewDecoder(request.Body).Decode(&userToRegister)

		if err != nil {
			writer.WriteHeader(400)
			fmt.Fprint(writer, "Body is not valid")
		}

		err = validate.Struct(userToRegister)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				writer.WriteHeader(400)
				fmt.Printf("A validation error occured: %v", err)
				fmt.Fprint(writer, err)

			}
		}
		url := firebaseRestBasePath + accountPath + ":signUp?key=" + firebaseAPIKey
		jsonBody, _ := json.Marshal(userToRegister)
		body := bytes.NewBuffer(jsonBody)
		resp, err := http.Post(url, "application/json", body)

		if err != nil {
			fmt.Fprintf(writer, "Failed to create user: %v", err)
		}
		var userThatsRegistered UserThatsRegistered
		err = json.UnMarshal(resp, userThatsRegistered)
		if err != nil {
			log.Printf("Firebase response not parsable. Reason: %v", err.error)
			fmt.Fprint(writer, "User not created")
		}
		defer resp.Body.Close()
		responseData, err := json.Marshal(UserThatsRegistered)
		fmt.Printf("Created user: %s", userThatsRegistered.idToken)

		writer.Header().Set("Content-Type", "application/json")

		writer.Write(responseData)
	} else {
		fmt.Fprint(writer, "Request method not supported.")
	}
}
