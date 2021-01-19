package auth

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	// SSMType defines server-side session management (session stored in server memory)
	SSMType = "session"
	// SSMCookieName defines cookie name
	SSMCookieName = "SESSIONID"
)

// DefineSSM performs Server-Side Session Management struct declaration
func DefineSSM(sessionDuration, pbkdf2Iterations, pbkdf2KeyLenght int) *SSM {
	return &SSM{
		sessionDuration:  sessionDuration,
		additionMutex:    sync.Mutex{},
		sessions:         make(map[string]Session),
		pbkdf2Iterations: pbkdf2Iterations,
		pbkdf2KeyLenght:  pbkdf2KeyLenght,
	}
}

// SSM is a server-side session management
// It will create a new session in server memory and puts
// session ID as a cookie
type SSM struct {
	// Sets session duration in seconds
	sessionDuration int
	// Session addition mutex
	additionMutex sync.Mutex
	// Map holds sessions
	sessions map[string]Session
	// Number of iterations that will be given to a
	// PBKDF2 hashing algorithm
	pbkdf2Iterations int
	// Resulted length of PBKDF2 key
	pbkdf2KeyLenght int
}

// CreateSession creates session and returns a session ID
func (ssm *SSM) CreateSession(w http.ResponseWriter) (string, error) {
	ssm.additionMutex.Lock()
	defer ssm.additionMutex.Unlock()
	sessionID, err := GenerateRandomString(16)
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %v", err)
	}
	for {
		if _, exist := ssm.sessions[sessionID]; exist {
			sessionID, err = GenerateRandomString(16)
			if err != nil {
				return "", fmt.Errorf("failed to generate session ID: %v", err)
			}
		} else {
			break
		}
	}

	cookie := http.Cookie{
		Name:     SSMCookieName,
		Value:    sessionID,
		MaxAge:   0,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
	ssm.sessions[sessionID] = Session{ID: sessionID, LastActivity: time.Now()}

	return sessionID, nil
}

// CheckSession checks if session is active or valid
// performs check on a cookie
func (ssm *SSM) CheckSession(w http.ResponseWriter, r *http.Request) error {
	if r == nil {
		return fmt.Errorf("request is nil")
	}

	cookie, err := r.Cookie(SSMCookieName)
	if err != nil {
		return fmt.Errorf("failed to get '%s' cookie: %v", SSMCookieName, err)
	}
	session, exist := ssm.sessions[cookie.Value]
	if !exist {
		return fmt.Errorf("session with '%s' ID does not exist", cookie.Value)
	}

	duration := time.Since(session.LastActivity)
	if duration.Seconds() > float64(ssm.sessionDuration) {
		delete(ssm.sessions, cookie.Value)
		return fmt.Errorf("session with '%s' ID has expired", cookie.Value)
	}
	// Update session last activity
	session.LastActivity = time.Now()
	ssm.sessions[cookie.Value] = session

	return nil
}

// Logout marks user session as ended
func (ssm *SSM) Logout(r *http.Request) error {
	if r == nil {
		return fmt.Errorf("request is nil")
	}

	cookie, err := r.Cookie(SSMCookieName)
	if err != nil {
		return fmt.Errorf("failed to get '%s' cookie: %v", SSMCookieName, err)
	}

	delete(ssm.sessions, cookie.Value)
	return nil
}

// PBKDF2HashPassword performs password hashing
func (ssm *SSM) PBKDF2HashPassword(password string, salt string) string {
	return pbkdf2HashPassword(password, salt, ssm.pbkdf2Iterations, ssm.pbkdf2KeyLenght)
}

// Session that stored in server memory
type Session struct {
	ID           string    `json:"id"`
	LastActivity time.Time `json:"lastActivity"`
}
