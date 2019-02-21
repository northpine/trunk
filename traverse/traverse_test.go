package traverse_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/northpine/trunk/traverse"
)

func TestEsriTraverser_GetLayers(t *testing.T) {
	client := http.Client{}
	URL, _ := url.Parse("https://gismaps.kingcounty.gov/arcgis/rest/services?f=json")
	traverser := traverse.NewEsriTraverser()
	layers, err := traverser.GetLayers(client, *URL, traverse.Esri)
	if err != nil {
		t.Fatal(err)
	}
	for _, layer := range layers {
		fmt.Println(layer.Name, " ", layer.URL)
	}
}
