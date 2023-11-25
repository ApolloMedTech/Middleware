package sessionmanager

import (
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
func (*MySessionStore) ReadState(*http.Request) (authboss.ClientState, error) {
	panic("unimplemented")
}

// WriteState implements authboss.ClientStateReadWriter.
func (*MySessionStore) WriteState(http.ResponseWriter, authboss.ClientState, []authboss.ClientStateEvent) error {
	panic("unimplemented")
}

// NewMySessionStore creates a new instance of MySessionStore.
func NewMySessionStore() *MySessionStore {
	return &MySessionStore{
		store: sessions.NewCookieStore([]byte("your-secret-key")),
	}
}
