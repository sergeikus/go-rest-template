package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LogInRequest represents a login request
type LogInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate performs LogInRequest fields validation
func (lor *LogInRequest) Validate() error {
	if len(lor.Username) == 0 {
		return fmt.Errorf("username must be provided")
	}
	if len(lor.Password) == 0 {
		return fmt.Errorf("password must be provided")
	}
	return nil
}

const logInTag = "LogIn"

// LogIn performs user log in
func (api *API) LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var lir LogInRequest
		if err := decoder.Decode(&lir); err != nil {
			fail(w, logInTag, fmt.Errorf("failed to decode body request: %v", err), http.StatusBadRequest)
			return
		}
		if err := lir.Validate(); err != nil {
			fail(w, logInTag, fmt.Errorf("log in request validation failed: %v", err), http.StatusBadRequest)
			return
		}

		passwordSalt, err := api.DB.GetUserSalt(lir.Username)
		if err != nil {
			fail(w, logInTag, fmt.Errorf("failed to get user salt: %v", err), http.StatusUnauthorized)
			return
		}

		passwordHash := api.Auth.PBKDF2HashPassword(lir.Password, passwordSalt)

		_, err = api.DB.VerifyUserCredentials(lir.Username, passwordHash)
		if err != nil {
			fail(w, logInTag, err, http.StatusUnauthorized)
			return
		}

		if _, err := api.Auth.CreateSession(w); err != nil {
			fail(w, logInTag, fmt.Errorf("failed to create session: %v", err), http.StatusUnauthorized)
			return
		}

		writeResponseString(w, MsgStatusOK, logInTag, fmt.Sprintf("Successfully logged in user with '%s' username", lir.Username))
	}
}

const logInStatusTag = "LogInStatus"

// LogInStatus checks if user is logged in or is authorized
func (api *API) LogInStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := api.Auth.CheckSession(w, r); err != nil {
			fail(w, logInStatusTag, err, http.StatusUnauthorized)
			return
		}

		writeResponseString(w, MsgStatusOK, logInStatusTag, fmt.Sprint("Session is active"))
	}
}

const logoutTag = "Logout"

// Logout performs logging out in server
func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := api.Auth.Logout(r); err != nil {
			fail(w, logoutTag, err, http.StatusInternalServerError)
			return
		}

		writeResponseString(w, MsgStatusOK, logoutTag, "Successfully logged out")
	}
}
