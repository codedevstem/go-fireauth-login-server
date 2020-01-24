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

// InUser is the struct that is expected when a user register for the application
type InUser struct {
	Email    string `json:"email" validate:"required,email,max=60"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
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

		inUser := InUser{}

		err := json.NewDecoder(request.Body).Decode(&inUser)

		if err != nil {
			writer.WriteHeader(400)
			fmt.Fprint(writer, "Body is not valid")
		}

		err = validate.Struct(inUser)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				writer.WriteHeader(400)
				fmt.Printf("A validation error occured: %v", err)
				fmt.Fprint(writer, err)

			}
		}
		url := firebaseRestBasePath + accountPath + ":signUp?key=" + firebaseAPIKey
		jsonBody, _ := json.Marshal(inUser)
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

		fmt.Printf("Created user: %s", inUser.Username)

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(responseData)
	} else {
		fmt.Fprint(writer, "Request method not supported.")
	}
}
