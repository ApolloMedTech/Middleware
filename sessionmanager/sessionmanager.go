package sessionmanager

import (
	"encoding/json"
	"fmt"
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

	// Use a closure to capture the request
	doWriteState := func(r *http.Request) error {
		// Serialize the ClientState to a JSON string
		serializedState, err := json.Marshal(state)
		if err != nil {
			return err
		}

		// Get the session from the request
		session, err := s.store.Get(r, "authsession")
		if err != nil {
			return err
		}

		// Store the serialized state in the session
		session.Values["authboss_state"] = string(serializedState)

		// Save the session
		err = session.Save(r, w)
		if err != nil {
			return err
		}

		return nil
	}

	// Perform the write state operation using the closure
	return doWriteState(nil) // Pass the actual request when calling this method
}

// NewMySessionStore creates a new instance of MySessionStore.
func NewMySessionStore() *MySessionStore {
	return &MySessionStore{
		store: sessions.NewCookieStore([]byte("your-secret-key")),
	}
}
