package auth

import (
	"database/sql"

	"github.com/ApolloMedTech/Middleware/config"
	"github.com/ApolloMedTech/Middleware/dbmanager"
	"github.com/volatiletech/authboss/v3"
	"golang.org/x/crypto/bcrypt"
)

func Login(email, password, userType string) (*config.ApolloUser, error) {

	dbManager, err := dbmanager.NewDBManager()
	if err != nil {
		return nil, err
	}

	defer dbManager.DB.Close()

	row := dbManager.DB.QueryRow("SELECT user_id, name, email, user_type FROM users WHERE email = $1 and password_hash = $2 and user_type = $3;", email, password, userType)

	var user config.ApolloUser

	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.UserType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, authboss.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
