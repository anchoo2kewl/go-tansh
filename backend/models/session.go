package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/anchoo2kewl/tansh.us/rand"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
	CreatedAt string
}

type SessionService struct {
	DB *sql.DB
}

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (ss *SessionService) Create(userID int) (*Session, error) {
	// TODO: Create the session token
	token, err := rand.SessionToken()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	hashedToken, err := ss.GenerateHashedToken(token)

	if err != nil {
		return nil, fmt.Errorf("create token, generate hashed token: %w", err)
	}

	id, err := ss.GetSession(userID)

	if err != nil {
		return nil, fmt.Errorf("create token, fetch existing token: %w", err)
	}

	row := ss.DB.QueryRow("")

	if id != 0 {
		fmt.Println("Session exists. Updating ...")
		row = ss.DB.QueryRow(`
		UPDATE sessions set token_hash = $1 where ID = $2 
		RETURNING id`, hashedToken, id)
	} else {
		fmt.Println("Session does not exist. Creating ...")
		row = ss.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2) RETURNING id`, userID, hashedToken)
	}

	session := Session{
		Token:     token,
		TokenHash: hashedToken,
	}

	err = row.Scan(&session.ID)

	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) GetSession(userID int) (int, error) {
	var session Session

	query := `SELECT * FROM sessions WHERE user_id = $1`
	row := ss.DB.QueryRow(query, userID)
	err := row.Scan(&session.ID, &session.UserID, &session.TokenHash, &session.CreatedAt)

	if err != nil {
		fmt.Errorf("User %d does not exist. Cannot create token: %w", userID, err)
	}
	return session.ID, nil
}

func (ss *SessionService) User(token string, email string) (*User, error) {

	email = strings.ToLower(email)

	user := User{
		Email: email,
	}

	session := Session{
		Token: token,
	}

	row := ss.DB.QueryRow(`SELECT s.id, s.token_hash from users as u INNER JOIN sessions as s
							on u.id = s.user_id where u.email like $1`, email)

	err := row.Scan(&session.ID, &session.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("session, email incorrect: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(session.TokenHash), []byte(token))
	if err != nil {
		return nil, fmt.Errorf("authenticate session: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Logout(email string) {

	email = strings.ToLower(email)

	ss.DB.QueryRow(`DELETE FROM sessions WHERE user_id in (select id from users where email like $1)`, email)

}

func (ss *SessionService) GenerateHashedToken(token string) (string, error) {
	hashedTokenBytes, err := bcrypt.GenerateFromPassword(
		[]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("create token: %w", err)
	}

	return string(hashedTokenBytes), nil
}
