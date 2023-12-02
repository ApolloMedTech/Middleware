package dbmanager

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ApolloMedTech/Middleware/config"
	_ "github.com/lib/pq" // PostgresSQL driver
	"github.com/sirupsen/logrus"
	authboss "github.com/volatiletech/authboss/v3"
)

// DBManager holds the database connection pool.
type DBManager struct {
	DB *sql.DB
}

// NewDBManager creates a new DBManager.
func NewDBManager() (*DBManager, error) {
	cfg := config.GetConfig().Database
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Name, cfg.Host, cfg.Port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.Errorf("Error opening database: %v", err)
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		logrus.Errorf("Error connecting to database: %v", err)
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return &DBManager{DB: db}, nil
}

// Insert executes an insert query and returns the ID of the last inserted row.
func (manager *DBManager) Insert(query string, args ...interface{}) (int64, error) {
	result, err := manager.DB.Exec(query, args...)
	if err != nil {
		logrus.Errorf("Error executing insert query '%s': %v", query, err)
		return 0, fmt.Errorf("error executing insert query '%s': %v", query, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		logrus.Errorf("Error getting last insert ID for query '%s': %v", query, err)
		return 0, fmt.Errorf("error getting last insert ID for query '%s': %v", query, err)
	}

	return id, nil
}

// Update executes an update query and returns the number of affected rows.
func (manager *DBManager) Update(query string, args ...interface{}) (int64, error) {
	result, err := manager.DB.Exec(query, args...)
	if err != nil {
		logrus.Errorf("Error executing update query '%s': %v", query, err)
		return 0, fmt.Errorf("error executing update query '%s': %v", query, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("Error getting rows affected for query '%s': %v", query, err)
		return 0, fmt.Errorf("error getting rows affected for query '%s': %v", query, err)
	}

	return rowsAffected, nil
}

// Delete executes a delete query and returns the number of affected rows.
func (manager *DBManager) Delete(query string, args ...interface{}) (int64, error) {
	result, err := manager.DB.Exec(query, args...)
	if err != nil {
		logrus.Errorf("Error executing delete query '%s': %v", query, err)
		return 0, fmt.Errorf("error executing delete query '%s': %v", query, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("Error getting rows affected for query '%s': %v", query, err)
		return 0, fmt.Errorf("error getting rows affected for query '%s': %v", query, err)
	}

	return rowsAffected, nil
}

// Select executes a select query and returns the rows.
func (manager *DBManager) Select(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := manager.DB.Query(query, args...)
	if err != nil {
		logrus.Errorf("Error executing select query '%s': %v", query, err)
		return nil, fmt.Errorf("error executing select query '%s': %v", query, err)
	}
	return rows, nil
}

// BeginTransaction starts a new database transaction.
func (manager *DBManager) BeginTransaction() (*sql.Tx, error) {
	tx, err := manager.DB.Begin()
	if err != nil {
		logrus.Errorf("Error beginning transaction: %v", err)
		return nil, fmt.Errorf("error beginning transaction: %v", err)
	}
	return tx, nil
}

// Close closes the database connection.
func (manager *DBManager) Close() error {
	if err := manager.DB.Close(); err != nil {
		logrus.Errorf("Error closing database connection: %v", err)
	}
	return nil
}

// AB - SERVER STORER

func (db *DBManager) Load(ctx context.Context, key string) (authboss.User, error) {
	// Use ConnectDB to establish a database connection
	dbManager, err := NewDBManager()
	if err != nil {
		return nil, err
	}

	defer dbManager.DB.Close()

	// Prepare SQL query to fetch user data based on username
	row := dbManager.DB.QueryRow("SELECT user_id, email FROM users WHERE email = $1;", key)

	var user config.ApolloUser
	// Scan the retrieved row into the User struct
	err = row.Scan(&user.ID, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, authboss.ErrUserNotFound // User not found
		}
		return nil, err // Other error while scanning
	}

	// Return the user object retrieved from the database
	return &user, nil
}

func (*DBManager) Save(ctx context.Context, user authboss.User) error {

	// Use ConnectDB to establish a database connection
	dbManager, err := NewDBManager()
	if err != nil {
		return err
	}
	defer dbManager.DB.Close()

	// Prepare SQL query for user insertion
	_, err = dbManager.DB.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2);",
		user.(*config.ApolloUser).Email, user.(*config.ApolloUser).Password)
	if err != nil {
		return err
	}

	// If the insertion was successful, return nil indicating no error
	return nil
}

// AB - RECOVERING SERVER STORER

func (*DBManager) LoadByRecoverSelector(ctx context.Context, selector string) (authboss.RecoverableUser, error) {
	// Use ConnectDB to establish a database connection
	dbManager, err := NewDBManager()
	if err != nil {
		return nil, err
	}

	defer dbManager.DB.Close()

	// Prepare SQL query to fetch user data based on selector
	row := dbManager.DB.QueryRow("SELECT id, email, recoverselector, recoververifier, recoveryexpiry FROM users WHERE recoverselector = $1;", selector)

	var user config.ApolloUser
	// Scan the retrieved row into the User struct
	err = row.Scan(&user.ID, &user.Email, &user.RecoverSelector, &user.RecoverVerifier, &user.RecoverExpiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, authboss.ErrUserNotFound // User not found
		}
		return nil, err // Other error while scanning
	}

	// Return the user object retrieved from the database
	return &user, nil
}
