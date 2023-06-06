package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/AkifhanIlgaz/lenslocked/rand"
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	token, err := rand.SessionToken()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	session := Session{
		UserID: userID,
		Token:  token,
	}

	hash := hmac.New(sha256.New, []byte(token))
	tokenHash := string(hash.Sum(nil))
	session.TokenHash = tokenHash

	row := ss.DB.QueryRow(`
		INSERT INTO sessions(user_id, token_hash)
		VALUES (
			$1,
			$2
		) RETURNING id;
	`, userID, tokenHash)

	err = row.Scan(&session.ID)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	//  TODO: Implement SessionService.User
	return nil, nil
}
