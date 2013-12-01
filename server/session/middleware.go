package session

import (
	"net/http"
)

// SessionHandler is a http.Handler that can also receive a username
type SessionHandler interface {
	http.Handler
	SetUsername(string)
}

// EnsureHasSession is a middleware that verifies every request
func EnsureHasSession(sessionMngr *SessionManager, tokenName String, next SessionHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		tokens := req.URL.Query()[tokenName]
		if len(tokens) > 1 {
			http.Error(rw, "Too many tokens.", http.StatusUnauthorized)
			return
		} else if len(tokens) == 0 {
			http.Error(rw, "No token.", http.StatusUnauthorized)
			return
		}
		token := tokens[0]

		username, ok := sessionMngr.UseSession(token, req.RemoteAddr)
		if !ok {
			http.Error(rw, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		next.SetUsername(username)
		next.ServeHTTP(rw, req)
	}
}
