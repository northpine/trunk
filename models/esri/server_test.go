package esri_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/northpine/trunk/models/esri"
)

func TestRoot(t *testing.T) {
	base, _ := url.Parse("https://gismaps.kingcounty.gov/arcgis/rest/services")
	root := esri.Root{}
	root.URL = *base
	b, _ := ioutil.ReadFile("data/root.json")
	err := json.Unmarshal(b, &root)
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{
		"/arcgis/rest/services/Property/KingCo_FarmlandPreservationProgram/MapServer",
		"/arcgis/rest/services/Property/KingCo_Parcels/MapServer",
		"/arcgis/rest/services/Property/KingCo_PropertyInfo/MapServer",
	}
	expected := make([]url.URL, len(paths))
	for i, p := range paths {
		expected[i] = url.URL{
			Scheme:   "https",
			Host:     "gismaps.kingcounty.gov",
			RawQuery: esri.FormatJSON,
			Path:     p,
		}
	}
	actual := root.GetChildren()
	for _, e := range expected {
		foundUrl := false
		for _, a := range actual {
			fmt.Println(a.String())
			if a.String() == e.String() {
				foundUrl = true
				break
			}
		}
		if !foundUrl {
			t.Fatalf("Could not find: %s", e.String())
		}
	}
}

func GetUrls(ss []string) []url.URL {
	urls := make([]url.URL, len(ss))
	for _, s := range ss {
		u, _ := url.Parse(s)
		urls = append(urls, *u)
	}
	return urls
}
