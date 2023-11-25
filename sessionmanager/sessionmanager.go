package sessionmanager

import (
	"fmt"
	"net/http"
	"strings"

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

// SessionState is an authboss.ClientState implementation that
// holds the request's session values for the duration of the request.
type SessionState struct {
	session *sessions.Session
}

// Get a key from the session
func (s SessionState) Get(key string) (string, bool) {
	str, ok := s.session.Values[key]
	if !ok {
		return "", false
	}
	value := str.(string)

	return value, ok
}
func (s MySessionStore) WriteState(w http.ResponseWriter, state authboss.ClientState, ev []authboss.ClientStateEvent) error {
	ses := state.(*SessionState)

	for _, ev := range ev {
		switch ev.Kind {
		case authboss.ClientStateEventPut:
			ses.session.Values[ev.Key] = ev.Value
		case authboss.ClientStateEventDel:
			delete(ses.session.Values, ev.Key)
		case authboss.ClientStateEventDelAll:
			if len(ev.Key) == 0 {
				ses.session.Options.MaxAge = -1
			} else {
				whitelist := strings.Split(ev.Key, ",")
				s.DeleteSessionValues(ses, whitelist)
			}
		}
	}

	return s.store.Save(nil, w, ses.session)
}

func (MySessionStore) DeleteSessionValues(ses *SessionState, whitelist []string) {
	for key := range ses.session.Values {
		if k, ok := key.(string); ok && !contains(whitelist, k) {
			delete(ses.session.Values, key)
		}
	}
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// NewMySessionStore creates a new instance of MySessionStore.
func NewMySessionStore() *MySessionStore {
	return &MySessionStore{
		store: sessions.NewCookieStore([]byte("your-secret-key")),
	}
}
