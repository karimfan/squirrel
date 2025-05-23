package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Store interface {
	CreateAccount(name, email, password, squirrel string) (int, error)
	Login(email, password string) (int, error)
	AddItem(accountID int, content string, tags []string, itemType string) (int, error)
}

type Server struct {
	store Store
}

func NewServer(s Store) *Server { return &Server{store: s} }

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", s.handleGraphQL)
	return mux
}

type gqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type gqlResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (s *Server) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	var req gqlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErr(w, err)
		return
	}
	query := req.Query
	switch {
	case strings.Contains(query, "createAccount"):
		s.createAccount(w, req.Variables)
	case strings.Contains(query, "login"):
		s.login(w, req.Variables)
	case strings.Contains(query, "addItem"):
		s.addItem(w, req.Variables)
	default:
		s.writeErr(w, fmt.Errorf("unknown operation"))
	}
}

func (s *Server) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s *Server) writeErr(w http.ResponseWriter, err error) {
	s.writeJSON(w, gqlResponse{Error: err.Error()})
}

func (s *Server) createAccount(w http.ResponseWriter, vars map[string]interface{}) {
	name := vars["name"].(string)
	email := vars["email"].(string)
	password := vars["password"].(string)
	squirrel := vars["squirrel"].(string)
	id, err := s.store.CreateAccount(name, email, password, squirrel)
	if err != nil {
		s.writeErr(w, err)
		return
	}
	s.writeJSON(w, gqlResponse{Data: map[string]interface{}{"createAccount": map[string]int{"id": id}}})
}

func (s *Server) login(w http.ResponseWriter, vars map[string]interface{}) {
	email := vars["email"].(string)
	password := vars["password"].(string)
	id, err := s.store.Login(email, password)
	if err != nil {
		s.writeErr(w, fmt.Errorf("invalid credentials"))
		return
	}
	token := fmt.Sprintf("token-%d", id)
	s.writeJSON(w, gqlResponse{Data: map[string]string{"login": token}})
}

func (s *Server) addItem(w http.ResponseWriter, vars map[string]interface{}) {
	token := vars["token"].(string)
	content := vars["content"].(string)
	var tags []string
	if v, ok := vars["tags"]; ok && v != nil {
		for _, t := range v.([]interface{}) {
			tags = append(tags, t.(string))
		}
	}
	var accountID int
	if _, err := fmt.Sscanf(token, "token-%d", &accountID); err != nil {
		s.writeErr(w, fmt.Errorf("invalid token"))
		return
	}
	itemType := "note"
	if strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://") {
		itemType = "article"
	}
	id, err := s.store.AddItem(accountID, content, tags, itemType)
	if err != nil {
		s.writeErr(w, err)
		return
	}
	s.writeJSON(w, gqlResponse{Data: map[string]interface{}{"addItem": map[string]int{"id": id}}})
}
