package login

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

const firebaseRestLoginPath = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"

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

type UserIn struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

func Handler(writer http.ResponseWriter, request *http.Request) {
	userIn := UserIn{}

	err := json.NewDecoder(request.Body).Decode(&userIn)
	if err != nil {
		fmt.Fprint(writer, "Body is not valid")
	}

	err = validate.Struct(userIn)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			writer.WriteHeader(400)
			fmt.Printf("A validation error occured: %v", err)
			fmt.Fprint(writer, err)
		}
	}

}
