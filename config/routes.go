package config

import "net/http"

func SetupRoutes() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/users", usersHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
