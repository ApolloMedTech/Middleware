package cookiemanager

import (
	"encoding/json"
	"fmt"
	"net/http"

	authboss "github.com/volatiletech/authboss/v3"
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

// ReadState implements authboss.ClientStateReadWriter.
func (m *MyCookieStore) ReadState(r *http.Request) (authboss.ClientState, error) {
	var state authboss.ClientState
	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		// Return an empty state if the cookie is not found (no error for missing cookie)
		if err == http.ErrNoCookie {
			return state, nil
		}
		return nil, err
	}

	// Decode the cookie value (assuming it's JSON)
	err = json.Unmarshal([]byte(cookie.Value), &state)
	if err != nil {
		return nil, fmt.Errorf("failed to decode client state cookie: %v", err)
	}

	return state, nil
}

// WriteState implements authboss.ClientStateReadWriter.
func (m *MyCookieStore) WriteState(w http.ResponseWriter, state authboss.ClientState, events []authboss.ClientStateEvent) error {
	// Encode the state as JSON
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to encode client state as JSON: %v", err)
	}

	// Create a new cookie with the encoded state
	cookie := http.Cookie{
		Name:  m.cookieName,
		Value: string(stateJSON),
		// Add other cookie options as needed
	}

	// Set the cookie in the response
	http.SetCookie(w, &cookie)

	return nil
}
