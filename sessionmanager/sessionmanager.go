package sessionmanager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ApolloMedTech/Middleware/dbmanager"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	authboss "github.com/volatiletech/authboss/v3"
)

// MySessionStore is a custom session store that implements the authboss.SessionState interface.
type MySessionStore struct {
	store sessions.Store
}

// SessionState is an authboss.ClientState implementation that
// holds the request's session values for the duration of the request.
type SessionState struct {
	session *sessions.Session
}

// NewMySessionStore creates a new instance of MySessionStore.
func NewMySessionStore() *MySessionStore {
	return &MySessionStore{
		store: sessions.NewCookieStore([]byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")),
	}
}

func (m *MySessionStore) CreateSession(userID int) (uuid.UUID, error) {

	// Use ConnectDB to establish a database connection
	dbManager, err := dbmanager.NewDBManager()
	if err != nil {
		return uuid.Nil, err
	}
	defer dbManager.DB.Close()

	token := uuid.New()

	// Prepare SQL query for user insertion
	_, err = dbManager.DB.Exec("INSERT INTO session (session_id, user_id, start_date, expiration_date) VALUES ($1, $2, $3, $4);", token, userID, time.Now(),
		time.Now().AddDate(0, 0, 1)) // fica por agora com um dia de sessão.
	if err != nil {
		return uuid.Nil, err
	}

	// If the insertion was successful, return nil indicating no error
	return token, nil
}

func (m *MySessionStore) InvalidateSession(sessionID uuid.UUID) error {

	// Use ConnectDB to establish a database connection
	dbManager, err := dbmanager.NewDBManager()
	if err != nil {
		return err
	}
	defer dbManager.DB.Close()

	// Prepare SQL query for user insertion
	_, err = dbManager.DB.Exec("UPDATE session SET active = 0 where session_id = $1;", sessionID) // fica por agora com um dia de sessão.
	if err != nil {
		return err
	}

	// If the insertion was successful, return nil indicating no error
	return nil
}

func (m *MySessionStore) IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {

	ssk, err := m.Load(w, r, "Session")

	if err != nil {
		logrus.Errorf("Error making request to microservice: %v", err)
		return false
	}

	if ssk == "" {
		logrus.Errorf("Sem sessão: %v", err)
		return false
	}

	return true
}

// func (m *MySessionStore) IsSessionStilValid(sessionID uuid.UUID) (bool, error) {
// 	// Use ConnectDB to establish a database connection
// 	dbManager, err := dbmanager.NewDBManager()
// 	if err != nil {
// 		return false, err
// 	}
// 	defer dbManager.DB.Close()

// 	// Prepare SQL query for user insertion
// 	_, err = dbManager.DB.Exec("select session_id from session where expiration_date <=  CURRENT_TIMESTAMP() and active = 1")
// 	if err != nil {
// 		return false, err
// 	}

// 	// If the insertion was successful, return nil indicating no error
// 	return true, nil
// }

// Save saves the session data for a given session token.
func (m *MySessionStore) Save(w http.ResponseWriter, r *http.Request, key, value string) error {
	session, err := m.store.Get(r, key)
	if err != nil {
		return err
	}

	// Store the user ID in the session, assuming it's stored as "user_id"
	session.Values[key] = value
	session.Options.MaxAge = 2 * 60 * 60

	// Save the session
	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (m *MySessionStore) DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := m.store.Get(r, "Session")
	if err != nil {
		return err
	}

	// Delete the session by setting its MaxAge to a negative value
	session.Options.MaxAge = -1
	err = session.Save(r, w)

	return err
}

func (m *MySessionStore) SaveObject(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value to JSON: %v", err)
	}

	m.Save(w, r, key, string(jsonValue))

	return nil
}

func (m *MySessionStore) Load(w http.ResponseWriter, r *http.Request, key string) (string, error) {
	session, err := m.store.Get(r, key)
	if err != nil {
		return "", err
	}

	// Extract the user ID from the session, assuming it's stored as "user_id"
	userID, ok := session.Values[key].(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in session")
	}

	return userID, nil
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
