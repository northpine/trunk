package insert

import (
	"io"
	"net/http"
)

// Servers inserts layers objects into our database
// Or queries for servers
func Servers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		Insert(r.Body)
	case "GET":
		QueryServers(w)
	}
}

//Insert inserts a server & its layers into a database
func Insert(r io.Reader) error {

	return nil
}

//QueryServers queries list of servers for matching terms
func QueryServers(w http.ResponseWriter) error {
	return nil
}
