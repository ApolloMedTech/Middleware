package cookiemanager

import (
	"fmt"
	"net/http"
)

// MyCookieStore is a custom cookie state store that implements the authboss.ClientState interface.
type MyCookieStore struct {
	cookieName string
}

// Load loads the client state data for a given client token from the cookie.
func (m *MyCookieStore) Load(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return "", fmt.Errorf("failed to load client state cookie: %v", err)
	}

	return cookie.Value, nil
}

// Save saves the client state data for a given client token to the cookie.
func (m *MyCookieStore) Save(w http.ResponseWriter, r *http.Request, token string) error {
	cookie := http.Cookie{
		Name:     m.cookieName,
		Value:    token,
		HttpOnly: true,
		// Add other cookie options as needed
	}

	http.SetCookie(w, &cookie)
	return nil
}

// NewMyCookieStore creates a new instance of MyCookieStore.
func NewMyCookieStore(cookieName string) *MyCookieStore {
	return &MyCookieStore{
		cookieName: cookieName,
	}
}

// AuthHandler is a placeholder handler for demonstration purposes.
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Handle authentication-related logic here
}
