package traverse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/northpine/trunk/errors"
	"github.com/northpine/trunk/models/esri"
)

// ServerType is the type of server being traversed
type ServerType string

const (
	//Esri server type
	Esri ServerType = "esri"
)

// EsriTraverser traverses an esri arcgis rest server
type EsriTraverser struct {
	layerLock *sync.Mutex
	layers    []*esri.Layer
	client    http.Client
	errors    errors.Errors
}

// NewEsriTraverser creates traverser
func NewEsriTraverser() *EsriTraverser {
	return &EsriTraverser{
		layerLock: &sync.Mutex{},
		errors:    errors.NewErrors(),
	}
}

// GetLayers traverses a server and returns layers
func (t *EsriTraverser) GetLayers(client http.Client, baseURL url.URL, serverType ServerType) ([]*esri.Layer, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%+v", r)
		}
	}()
	t.client = client
	t.Drill(baseURL)
	return t.layers, t.errors.Get()
}

//Drill all the way down an ESRI server
func (t *EsriTraverser) Drill(URL url.URL) {
	URL.RawQuery = esri.FormatJSON
	resp, err := t.client.Get(URL.String())
	if t.errors.Add(err) {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if t.errors.Add(err) {
		return
	}
	d, err := esri.GetDriller(URL, body)
	// Might be some kind of thing we don't support (like GPServer)
	if err != nil {
		return
	}
	if extracter, ok := esri.AsExtracter(d); ok {
		extracter.Extract(URL)
	}
	if layer, ok := esri.AsLayer(d); ok {
		layer.URL = URL
		t.addLayers(layer)
		return
	}
	if d.CanDrill() {
		children := d.GetChildren()
		for _, URL := range children {
			t.Drill(URL)
		}
	}
}

func (t *EsriTraverser) addLayers(layer *esri.Layer) {
	t.layerLock.Lock()
	defer t.layerLock.Unlock()
	t.layers = append(t.layers, layer)
}
