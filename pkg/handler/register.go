package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sergeikus/go-rest-template/pkg/auth"
	"github.com/sergeikus/go-rest-template/pkg/types"
)

// RegisterUserRequest is a user request for a registration
type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

var errRegisterUserRequestNoUsername = errors.New("username is not provided")
var errRegisterUserRequestNoPassword = errors.New("password is not provided")

// Validate performs register user request validation
func (rur *RegisterUserRequest) Validate() error {
	if len(rur.Username) == 0 {
		return errRegisterUserRequestNoUsername
	}
	if len(rur.Password) == 0 {
		return errRegisterUserRequestNoPassword
	}
	return nil
}

const registerUserTag = "RegisterUser"

// RegisterUser performs user registration in the database
func (api *API) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var rur RegisterUserRequest
		if err := decoder.Decode(&rur); err != nil {
			fail(w, registerUserTag, fmt.Errorf("failed to decode body request: %v", err), http.StatusBadRequest)
			return
		}
		if err := rur.Validate(); err != nil {
			fail(w, registerUserTag, fmt.Errorf("register user request validation failed: %v", err), http.StatusBadRequest)
			return
		}

		passwordSalt, err := auth.GenerateRandomString(16)
		if err != nil {
			fail(w, registerUserTag, fmt.Errorf("failed to generate password salt: %v", err), http.StatusInternalServerError)
			return
		}

		passwordHash := api.Auth.PBKDF2HashPassword(rur.Password, passwordSalt)

		user := types.User{
			Username:     rur.Username,
			PasswordSalt: passwordSalt,
			PasswordHash: passwordHash,
			Email:        rur.Email,
			IsDisabled:   false,
		}
		if _, err := api.DB.RegisterUser(user); err != nil {
			fail(w, logInTag, fmt.Errorf("failed to register new user: %v", err), http.StatusInternalServerError)
			return
		}

		writeResponseString(w, MsgStatusOK, registerUserTag, fmt.Sprintf("Successfully registered user with '%s' username", rur.Username))
	}
}
