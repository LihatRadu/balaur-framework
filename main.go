package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/yuin/gopher-lua"
)

//go:embed lua/*.lua
var luaScripts embed.FS


func executeLuaScript(script string) (string, error) {
    L := lua.NewState()
    defer L.Close()

    if err := L.DoString(script); err != nil {
        return "", err
    }

    // Assume the Lua script returns a string
    return L.Get(-1).(lua.LString).String(), nil
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        
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
