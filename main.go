package main

import (
    "fmt"
    "net/http"
    "github.com/yuin/gopher-lua"
)

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
        luaScript := `
            return "Hello from Lua!"
        `
        result, err := executeLuaScript(luaScript)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, result)
    })

    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
