package server

import "strings"

type Account struct {
	ID       int
	Name     string
	Email    string
	Password string
	Squirrel string
}

type Item struct {
	ID        int
	AccountID int
	Content   string
	Tags      string
	Type      string
}

type MemoryStore struct {
	Accounts      []Account
	Items         []Item
	nextAccountID int
	nextItemID    int
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{nextAccountID: 1, nextItemID: 1} }

func (m *MemoryStore) CreateAccount(name, email, password, squirrel string) (int, error) {
	id := m.nextAccountID
	m.nextAccountID++
	m.Accounts = append(m.Accounts, Account{ID: id, Name: name, Email: email, Password: password, Squirrel: squirrel})
	return id, nil
}

func (m *MemoryStore) Login(email, password string) (int, error) {
	for _, a := range m.Accounts {
		if a.Email == email && a.Password == password {
			return a.ID, nil
		}
	}
	return 0, sqlErr("invalid credentials")
}

type sqlErr string

func (e sqlErr) Error() string { return string(e) }

func (m *MemoryStore) AddItem(accountID int, content string, tags []string, itemType string) (int, error) {
	id := m.nextItemID
	m.nextItemID++
	m.Items = append(m.Items, Item{ID: id, AccountID: accountID, Content: content, Tags: strings.Join(tags, ","), Type: itemType})
	return id, nil
}
