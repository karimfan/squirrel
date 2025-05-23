package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/squirrel?sslmode=disable"
	}
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("ping: %v", err)
	}

	http.HandleFunc("/graphql", handleGraphQL)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type gqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type gqlResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	var req gqlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, err)
		return
	}

	query := req.Query
	switch {
	case strings.Contains(query, "createAccount"):
		createAccount(w, req.Variables)
	case strings.Contains(query, "login"):
		login(w, req.Variables)
	case strings.Contains(query, "addItem"):
		addItem(w, req.Variables)
	default:
		writeErr(w, fmt.Errorf("unknown operation"))
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	writeJSON(w, gqlResponse{Error: err.Error()})
}

func createAccount(w http.ResponseWriter, vars map[string]interface{}) {
	name := vars["name"].(string)
	email := vars["email"].(string)
	password := vars["password"].(string)
	squirrel := vars["squirrel"].(string)

	var id int
	err := db.QueryRow(`INSERT INTO accounts (name,email,password,squirrel_name) VALUES ($1,$2,$3,$4) RETURNING id`,
		name, email, password, squirrel).Scan(&id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, gqlResponse{Data: map[string]interface{}{"createAccount": map[string]int{"id": id}}})
}

func login(w http.ResponseWriter, vars map[string]interface{}) {
	email := vars["email"].(string)
	password := vars["password"].(string)
	var id int
	err := db.QueryRow(`SELECT id FROM accounts WHERE email=$1 AND password=$2`, email, password).Scan(&id)
	if err != nil {
		writeErr(w, fmt.Errorf("invalid credentials"))
		return
	}
	token := fmt.Sprintf("token-%d", id)
	writeJSON(w, gqlResponse{Data: map[string]string{"login": token}})
}

func addItem(w http.ResponseWriter, vars map[string]interface{}) {
	token := vars["token"].(string)
	content := vars["content"].(string)
	var tags []string
	if v, ok := vars["tags"]; ok {
		for _, t := range v.([]interface{}) {
			tags = append(tags, t.(string))
		}
	}
	var accountID int
	if _, err := fmt.Sscanf(token, "token-%d", &accountID); err != nil {
		writeErr(w, fmt.Errorf("invalid token"))
		return
	}
	itemType := "note"
	if strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://") {
		itemType = "article"
	}
	var id int
	tagString := strings.Join(tags, ",")
	err := db.QueryRow(`INSERT INTO items (account_id,content,tags,type) VALUES ($1,$2,$3,$4) RETURNING id`,
		accountID, content, tagString, itemType).Scan(&id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, gqlResponse{Data: map[string]interface{}{"addItem": map[string]int{"id": id}}})
}
