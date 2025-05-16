package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	tmpl      *template.Template
	tmplShell *template.Template
)

type PageData struct {
	Command string
	Stdout  string
	Stderr  string
	Dir     string
}

func init() {
	var err error
	tmpl, err = template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	tmplShell, err = template.ParseFiles(filepath.Join("templates", "shell.html"))
	if err != nil {
		log.Fatalf("Error parsing shell template: %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	data := PageData{
		Command: "",
		Stdout:  "",
		Stderr:  "",
		Dir:     dir,
	}
	tmpl.Execute(w, data)
}

func shellHandler(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	data := PageData{
		Command: "",
		Stdout:  "",
		Stderr:  "",
		Dir:     dir,
	}
	if err := tmplShell.Execute(w, data); err != nil {
		http.Error(w, "Error loading shell template", http.StatusInternalServerError)
	}
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	command := strings.TrimSpace(r.FormValue("command"))
	dir, _ := os.Getwd()

	if command == "" {
		data := PageData{
			Command: "",
			Stdout:  "",
			Stderr:  "",
			Dir:     dir,
		}
		if err := tmplShell.Execute(w, data); err != nil {
			http.Error(w, "Error loading shell template", http.StatusInternalServerError)
		}
		return
	}

	if strings.HasPrefix(command, "cd ") {
		path := strings.TrimSpace(command[3:])
		if strings.HasPrefix(path, "~") {
			path = os.Getenv("HOME")
			path += command[4:]
			os.Chdir(path)
		}
		if err := os.Chdir(path); err != nil {
			data := PageData{
				Command: command,
				Stdout:  "",
				Stderr:  err.Error(),
				Dir:     dir,
			}
			if err := tmplShell.Execute(w, data); err != nil {
				http.Error(w, "Error loading shell template", http.StatusInternalServerError)
			}
			return
		}
		dir, _ = os.Getwd()
		data := PageData{
			Command: command,
			Stdout:  "Changed directory to " + dir,
			Stderr:  "",
			Dir:     dir,
		}
		if err := tmplShell.Execute(w, data); err != nil {
			http.Error(w, "Error loading shell template", http.StatusInternalServerError)
		}
		return
	}

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	stdout := string(output)
	stderr := ""
	if err != nil {
		stderr = err.Error()
	}

	data := PageData{
		Command: command,
		Stdout:  stdout,
		Stderr:  stderr,
		Dir:     dir,
	}
	if err := tmplShell.Execute(w, data); err != nil {
		http.Error(w, "Error loading shell template", http.StatusInternalServerError)
	}
}

func executeShellHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		shellHandler(w, r)
	} else if r.Method == http.MethodPost {
		executeHandler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/execute", executeShellHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
