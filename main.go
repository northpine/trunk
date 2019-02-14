package main

import (
	"net/http"

	"github.com/northpine/trunk/insert"
	"github.com/northpine/trunk/traverse"
)

func main() {
	http.HandleFunc("/servers", insert.Servers)
	http.HandleFunc("/traverse", traverse.Traverse)
	http.ListenAndServe(":4000", nil)
}
