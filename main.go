package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	serverURL   = "http://localhost:8080/graphql"
	dataDir     = "data"
	sessionFile = filepath.Join(dataDir, "session_token")
)

type gqlResponse struct {
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
}

func sendMutation(query string, vars map[string]interface{}, out interface{}) error {
	payload := map[string]interface{}{"query": query, "variables": vars}
	b, _ := json.Marshal(payload)
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var r gqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Error != "" {
		return fmt.Errorf(r.Error)
	}
	if out != nil {
		return json.Unmarshal(r.Data, out)
	}
	return nil
}

func saveToken(token string) error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(sessionFile, []byte(token), 0644)
}

func loadToken() (string, error) {
	b, err := ioutil.ReadFile(sessionFile)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func runCreateAccount(args []string) error {
	fs := flag.NewFlagSet("create-account", flag.ExitOnError)
	name := fs.String("name", "", "name")
	email := fs.String("email", "", "email")
	password := fs.String("password", "", "password")
	squirrel := fs.String("squirrel", "", "squirrel name")
	fs.Parse(args)
	if *name == "" || *email == "" || *password == "" || *squirrel == "" {
		return fmt.Errorf("all fields required")
	}
	q := "mutation($name:String!,$email:String!,$password:String!,$squirrel:String!){createAccount(name:$name,email:$email,password:$password,squirrel:$squirrel){id}}"
	vars := map[string]interface{}{"name": *name, "email": *email, "password": *password, "squirrel": *squirrel}
	var resp struct{ CreateAccount struct{ ID int } }
	if err := sendMutation(q, vars, &resp); err != nil {
		return err
	}
	fmt.Printf("account created with id %d\n", resp.CreateAccount.ID)
	return nil
}

func runLogin(args []string) error {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	email := fs.String("email", "", "email")
	password := fs.String("password", "", "password")
	fs.Parse(args)
	if *email == "" || *password == "" {
		return fmt.Errorf("email and password required")
	}
	q := "mutation($email:String!,$password:String!){login(email:$email,password:$password)}"
	vars := map[string]interface{}{"email": *email, "password": *password}
	var resp struct{ Login string }
	if err := sendMutation(q, vars, &resp); err != nil {
		return err
	}
	if err := saveToken(resp.Login); err != nil {
		return err
	}
	fmt.Println("logged in")
	return nil
}

func runAddItem(args []string) error {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	tags := fs.String("tags", "", "comma separated tags")
	fs.Parse(args)
	if fs.NArg() < 1 {
		return fmt.Errorf("content required")
	}
	content := fs.Arg(0)
	token, err := loadToken()
	if err != nil {
		return fmt.Errorf("please log in first")
	}
	var tagList []string
	if *tags != "" {
		for _, t := range strings.Split(*tags, ",") {
			tt := strings.TrimSpace(t)
			if tt != "" {
				tagList = append(tagList, tt)
			}
		}
	}
	q := "mutation($token:String!,$content:String!,$tags:[String!]){addItem(token:$token,content:$content,tags:$tags){id}}"
	vars := map[string]interface{}{"token": token, "content": content, "tags": tagList}
	var resp struct{ AddItem struct{ ID int } }
	if err := sendMutation(q, vars, &resp); err != nil {
		return err
	}
	fmt.Printf("added item with id %d\n", resp.AddItem.ID)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("squirrel <command> [options]")
		fmt.Println("commands: create-account, login, add")
		return
	}
	var err error
	switch os.Args[1] {
	case "create-account":
		err = runCreateAccount(os.Args[2:])
	case "login":
		err = runLogin(os.Args[2:])
	case "add":
		err = runAddItem(os.Args[2:])
	default:
		fmt.Println("unknown command")
		return
	}
	if err != nil {
		fmt.Println("error:", err)
	}
}
