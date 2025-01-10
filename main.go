package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/yuin/gopher-lua"
)

//go:embed lua/*.lua
var luaScripts embed.FS

// Function to execute an embedded Lua script
func executeLuaScript(filename string) (string, error) {
	// Read the embedded Lua file
	data, err := luaScripts.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read Lua file: %v", err)
	}

	// Create a new Lua state
	L := lua.NewState()
	defer L.Close()

	// Execute the Lua script
	if err := L.DoString(string(data)); err != nil {
		return "", fmt.Errorf("failed to execute Lua script: %v", err)
	}

	// Assume the Lua script returns a string
	return L.Get(-1).(lua.LString).String(), nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Execute the embedded Lua script
		result, err := executeLuaScript("lua/balaur.lua")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, result)
	})

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
