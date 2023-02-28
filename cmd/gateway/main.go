package main

import (
    "flag"
	"encoding/json"
	"log"
	"net/http"
	"github.com/rs/cors"
	"github.com/apex/gateway"
    "fmt"
)

func main() {

	port := flag.Int("port", -1, "specify a port to use http rather than AWS Lambda")
    flag.Parse()
    listener := gateway.ListenAndServe
    portStr := "n/a"
    if *port != -1 {
        portStr = fmt.Sprintf(":%d", *port)
        listener = http.ListenAndServe
    }

	http.HandleFunc("/hello", createHttpHandler(handleHello))
	http.HandleFunc("/projects", createHttpHandler(handleProjects))
	http.HandleFunc("/users", createHttpHandler(handleUsers))
	http.HandleFunc("/commits", createHttpHandler(handleCommits))
	log.Fatal(listener(portStr, nil))
}

// Create a new HTTP handler using http.HandlerFunc and c.Handler
func createHttpHandler(handlerFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedHeaders: []string{"*"},
        AllowedMethods: []string{"*"},
    })
	handler := http.HandlerFunc(handlerFunc)
	handler = c.Handler(handler).(http.HandlerFunc)
	return handler
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode("Hello world!")
}

func handleProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Make an HTTP GET request to the projects endpoint
	owner := r.URL.Query().Get("owner")
	req, err := http.NewRequest("GET", "https://gitlab.com/api/v4/projects?owned="+owner, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Private-Token", "glpat-fXk2b-zcQNWsc9mjJsnb")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the response body to extract the list of projects
	var projects []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		HTTPURL  string `json:"http_url_to_repo"`
		Visibility string `json:"visibility"`
	}
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	projectId := r.URL.Query().Get("projectid")
    req, err := http.NewRequest("GET", "https://gitlab.com/api/v4/projects/" + projectId + "/members", nil)
    if err != nil {
        
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    req.Header.Add("Private-Token", "glpat-fXk2b-zcQNWsc9mjJsnb")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Parse the response body to extract the list of users
    var users []struct {
        ID       int    `json:"id"`
        Username string `json:"username"`
    }

    err = json.NewDecoder(resp.Body).Decode(&users)

	json.NewEncoder(w).Encode(users)
}

func handleCommits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	projectId := r.URL.Query().Get("projectid")
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
    req, err := http.NewRequest("GET", "https://gitlab.com/api/v4/projects/" + projectId + "/repository/commits?author_username=" + user, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    req.Header.Add("Private-Token", "glpat-fXk2b-zcQNWsc9mjJsnb")
	
    resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

    // Parse the response body to extract the list of commits
    var commits []struct {
        ID           string `json:"id"`
        Title        string `json:"title"`
        AuthorName   string `json:"author_name"`
        AuthorEmail  string `json:"author_email"`
        CommittedDate string `json:"committed_date"`
    }
    err = json.NewDecoder(resp.Body).Decode(&commits)

	json.NewEncoder(w).Encode(commits)
}