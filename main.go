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
)

//go:embed templates/*.html.lua
var templateFiles embed.FS

// TemplateService implements TemplateExecutor
type TemplateService struct{}

func (ts *TemplateService) ExecuteTemplate(filename string, data map[string]interface{}) (string, error) {
	content, err := templateFiles.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	registerLuaFunctions(L, data)
	output := parseTemplate(string(content), L)
	return output, nil
}

func registerLuaFunctions(L *lua.LState, data map[string]interface{}) {
	for key, value := range data {
		L.SetGlobal(key, toLuaValue(L, value))
	}

	L.SetGlobal("embed", L.NewFunction(func(L *lua.LState) int {
		expression := L.ToString(1)
		L.Push(lua.LString(fmt.Sprintf("%v", expression)))
		return 1
	}))
}

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

func parseTemplate(template string, L *lua.LState) string {
	re := regexp.MustCompile(`<%([\s\S]*?)%>`)
	output := re.ReplaceAllStringFunc(template, func(match string) string {
		code := strings.TrimSpace(match[2 : len(match)-2])
		if strings.HasPrefix(code, "=") {
			expression := strings.TrimSpace(code[1:])
			if err := L.DoString(fmt.Sprintf(`return embed(%q)`, expression)); err != nil {
				return fmt.Sprintf("<!-- ERROR: %v -->", err)
			}
			return L.Get(-1).String()
		} else {
			if err := L.DoString(code); err != nil {
				return fmt.Sprintf("<!-- ERROR: %v -->", err)
			}
			return ""
		}
	})
	return output
}

func main() {
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.DB.Close()

	if err := database.CreateUseerTable(config.DB); err != nil {
		log.Fatalf("Failed to create user table: %v", err)
	}

	// Initialize handlers with template executor
	config.InitHandlers(&TemplateService{})
	config.SetupRoutes()

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
