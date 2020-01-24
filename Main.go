package main
 
import (
	"net/http"

	"github.com/codedevstem/go-fireauth-login-server/register"
	"github.com/codedevstem/go-fireauth-login-server/login"
)

const port = ":13300"

// LoginUser is the struct that is expected when a user preforms a login
type LoginUser struct {
	Email    string `json:"Email" validate:"required,email,max=60"`
	Password string `json:"password" validate:"required,min=8"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", login.Handler)
	mux.HandleFunc("/register", register.Handler)
	mux.HandleFunc("/logout", logoutHandler)

	http.ListenAndServe(port, mux)

}
