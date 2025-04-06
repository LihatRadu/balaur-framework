package config

import (
	"fmt"
	"net/http"
	"balaur/database"
)

// TemplateExecutor defines the interface for executing templates
type TemplateExecutor interface {
	ExecuteTemplate(filename string, data map[string]interface{}) (string, error)
}

var (
	templateExec TemplateExecutor
)

// InitHandlers initializes the handlers with required dependencies
func InitHandlers(te TemplateExecutor) {
	templateExec = te
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"title":    "Balaur template",
		"greeting": "Balaur reveals!",
		"items":    []string{"Item 1", "Item 2", "Item 3"},
	}

	result, err := templateExec.ExecuteTemplate("templates/template.html.lua", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, result)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []database.User
	for rows.Next() {
		var user database.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	data := map[string]interface{}{
		"title": "User List",
		"users": users,
	}

	result, err := templateExec.ExecuteTemplate("templates/users.html.lua", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, result)
}
