package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/northpine/trunk/models/esri"

	"github.com/northpine/trunk/insert"
	"github.com/northpine/trunk/traverse"
)

func main() {
	http.HandleFunc("/servers", insert.Servers)
	http.HandleFunc("/traverse", traverse.Traverse)
	f, _ := os.Open("/Users/myles/output.csv")
	s := bufio.NewScanner(f)
	client := http.Client{
		Timeout: 1 * time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	dec := make(chan bool)
	num := 0
	servers := make([]string, 0)
	for s.Scan() {
		servers = append(servers, s.Text())
	}
	for _, server := range servers {
		num++
		go func(server string) {
			defer func() {
				dec <- true
			}()
			serverUrl, err := url.Parse(server)
			if err != nil {
				return
			}
			serverUrl.RawQuery = esri.FormatJSON
			traverser := traverse.NewEsriTraverser()
			layers, err := traverser.GetLayers(client, *serverUrl, traverse.Esri)
			if err != nil {
				return
			}
			if len(layers) == 0 {
				return
			}
			for _, layer := range layers {
				layer.URL.RawQuery = ""
				layer.DisplayURL = layer.URL.String()
			}
			insert := esri.Insert{
				Server: server,
				Layers: layers,
			}
			body, _ := json.Marshal(insert)

			buffer := bytes.NewBuffer(body)
			resp, err := http.Post("https://pine.center/insert", "application/json", buffer)
			if err != nil {
				log.Printf("Error in POST for %s: %+v", server, err)
			}
			respBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("%s : %s\n", server, string(respBody))
		}(server)
	}
	for i := 0; i < num; i++ {
		<-dec
	}

}
