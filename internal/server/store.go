package server

import (
	"database/sql"
	"strings"
)

type SQLStore struct {
	DB *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore { return &SQLStore{DB: db} }

func (s *SQLStore) CreateAccount(name, email, password, squirrel string) (int, error) {
	var id int
	err := s.DB.QueryRow(`INSERT INTO accounts (name,email,password,squirrel_name) VALUES ($1,$2,$3,$4) RETURNING id`,
		name, email, password, squirrel).Scan(&id)
	return id, err
}

func (s *SQLStore) Login(email, password string) (int, error) {
	var id int
	err := s.DB.QueryRow(`SELECT id FROM accounts WHERE email=$1 AND password=$2`, email, password).Scan(&id)
	return id, err
}

func (s *SQLStore) AddItem(accountID int, content string, tags []string, itemType string) (int, error) {
	var id int
	tagString := strings.Join(tags, ",")
	err := s.DB.QueryRow(`INSERT INTO items (account_id,content,tags,type) VALUES ($1,$2,$3,$4) RETURNING id`,
		accountID, content, tagString, itemType).Scan(&id)
	return id, err
}
