package main

import (
	"balaur/config"
	"balaur/database"
	"embed"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/yuin/gopher-lua"

//	"balaur-framework/config"
//	"balaur-framework/database"
)

//go:embed templates/*.html.lua
var templateFiles embed.FS

// Function to parse and execute a Lua-embedded HTML template
func executeTemplate(filename string, data map[string]interface{}) (string, error) {
	// Read the embedded template file
	content, err := templateFiles.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	// Create a new Lua state
	L := lua.NewState()
	defer L.Close()

	// Register Lua functions for embedding data
	registerLuaFunctions(L, data)

	// Parse the template
	output := parseTemplate(string(content), L)

	return output, nil
}

// Function to register Lua functions for embedding data
func registerLuaFunctions(L *lua.LState, data map[string]interface{}) {
	// Add data to the Lua environment
	for key, value := range data {
		L.SetGlobal(key, toLuaValue(L, value))
	}

	// Add a helper function for embedding expressions
	L.SetGlobal("embed", L.NewFunction(func(L *lua.LState) int {
		expression := L.ToString(1)
		L.Push(lua.LString(fmt.Sprintf("%v", expression)))
		return 1
	}))
}

// Function to convert Go values to Lua values
func toLuaValue(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case string:
		return lua.LString(v)
	case int:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case []interface{}:
		table := L.NewTable()
		for _, item := range v {
			table.Append(toLuaValue(L, item))
		}
		return table
	default:
		return lua.LNil
	}
}

// Function to parse the template and execute embedded Lua code
func parseTemplate(template string, L *lua.LState) string {
	// Regex to match Lua code blocks
	re := regexp.MustCompile(`<%([\s\S]*?)%>`)
	output := re.ReplaceAllStringFunc(template, func(match string) string {
		// Extract Lua code from the match
		code := strings.TrimSpace(match[2 : len(match)-2])

		// Execute Lua code
		if strings.HasPrefix(code, "=") {
			// Embed the result of the expression
			expression := strings.TrimSpace(code[1:])
			if err := L.DoString(fmt.Sprintf(`return embed(%q)`, expression)); err != nil {
				return fmt.Sprintf("<!-- ERROR: %v -->", err)
			}
			return L.Get(-1).String()
		} else {
			// Execute Lua code without embedding
			if err := L.DoString(code); err != nil {
				return fmt.Sprintf("<!-- ERROR: %v -->", err)
			}
			return ""
		}
	})

	return output
}

func homeHandler(w http.ResponseWriter, r *http.Request)  {
    if r.URL.Path != "/" {
        http.NotFound(w,r)
        return
    }

    data := map[string] interface{}{
        "title": "Balaur template",
        "greeting": "Balaur reveals!",
        "items":  []string{"Item 1", "Item 2", "Item 3"},
    }

    result, err := executeTemplate("templates/template.html.lua", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, result)
}

func usersHandler(w http.ResponseWriter, r *http.Request)  {
    rows, err := config.DB.Query("SELECT id, username, email FROM users")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []database.User

    for rows.Next(){
        var user database.User
        if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        users = append(users, user)
    }

// Data to pass to the template
	data := map[string]interface{}{
		"title": "User List",
		"users": users,
	}

	// Execute the template
	result, err := executeTemplate("templates/users.html.lua", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, result)
}

func main() {
    config.InitDB()
    defer config.DB.Close()

    if err := database.CreateUseerTable(&config.DB); err != nil {
        log.Fatalf("Failed to create user table: %v", err)
    }

    config.SetupRoutes()

    fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
