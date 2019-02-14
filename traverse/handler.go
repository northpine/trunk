package traverse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/northpine/trunk/models/esri"
)

// Traverse handles a request for a traverse of an ESRI server
func Traverse(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if servers, ok := params["server"]; ok && len(servers) > 0 {
		// Only need first one
		server := servers[0]
		fmt.Printf("Request traversal of %s\n", server)
		traverser := NewEsriTraverser()
		client := http.Client{
			Timeout: 1 * time.Minute,
		}
		serverURL, err := url.Parse(server)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%+v", err)
		}
		serverURL.RawQuery = esri.FormatJSON
		layers, err := traverser.GetLayers(client, *serverURL, Esri)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%+v", err)
		}
		for _, layer := range layers {
			layer.URL.RawQuery = ""
			layer.DisplayURL = layer.URL.String()
		}
		data := esri.Insert{
			Server: server,
			Layers: layers,
		}
		resp, _ := json.Marshal(data)
		w.Write(resp)
		fmt.Printf("Finished traversal of %s\n", server)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing server query param")
	}
}
