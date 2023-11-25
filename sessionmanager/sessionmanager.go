package sessionmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	authboss "github.com/volatiletech/authboss/v3"
)

// MySessionStore is a custom session store that implements the authboss.SessionState interface.
type MySessionStore struct {
	store sessions.Store
}

func (m *MySessionStore) Load(w http.ResponseWriter, r *http.Request, key string) (string, error) {
	session, err := m.store.Get(r, key)
	if err != nil {
		return "", err
	}

	// Extract the user ID from the session, assuming it's stored as "user_id"
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in session")
	}

	return userID, nil
}

// Save saves the session data for a given session token.
func (m *MySessionStore) Save(w http.ResponseWriter, r *http.Request, key, value string) error {
	session, err := m.store.Get(r, key)
	if err != nil {
		return err
	}

	// Store the user ID in the session, assuming it's stored as "user_id"
	session.Values["user_id"] = value

	// Save the session
	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

// ReadState implements authboss.ClientStateReadWriter.
func (m *MySessionStore) ReadState(r *http.Request) (authboss.ClientState, error) {
	// Retrieve the session from the store using the request
	var state authboss.ClientState

	session, err := m.store.Get(r, "my_session_name")
	if err != nil {
		// Return an empty state if the session is not found (no error for missing session)
		if err == http.ErrNoCookie {
			return state, nil
		}
		return nil, fmt.Errorf("failed to read client state from session: %v", err)
	}

	// Extract the client state from the session
	state, ok := session.Values["clientState"].(authboss.ClientState)
	if !ok {
		// Return an empty state if the client state is not found in the session
		return state, nil
	}

	return state, nil
}

// WriteState implements the authboss.SessionStorer interface
func (s MySessionStore) WriteState(w http.ResponseWriter, state authboss.ClientState, ev []authboss.ClientStateEvent) error {
	// Implement the logic to write the client state to the response writer

	// For example, you might serialize the state and write it to a cookie
	// or include it in the response headers.

	// Here is a simple example using a cookie:
	serializedState := serializeState(state) // Implement this function

	cookie := http.Cookie{
		Name:  "authboss_state_cookie",
		Value: serializedState,
		// Add other cookie options as needed (e.g., MaxAge, Path, etc.)
	}

	http.SetCookie(w, &cookie)

	return nil
}

// serializeState is a function that serializes authboss.ClientState into a JSON string
func serializeState(state authboss.ClientState) string {
	// Serialize the ClientState to a JSON string
	serializedState, err := json.Marshal(state)
	if err != nil {
		// Handle the error, for example, log it and return an empty string
		log.Printf("Error serializing client state: %v", err)
		return ""
	}

	return string(serializedState)
}

// NewMySessionStore creates a new instance of MySessionStore.
func NewMySessionStore() *MySessionStore {
	return &MySessionStore{
		store: sessions.NewCookieStore([]byte("your-secret-key")),
	}
}
