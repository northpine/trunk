package esri_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/northpine/trunk/models/esri"
)

func TestGetDriller(t *testing.T) {
	layer := esri.Layer{
		ID: 0,
	}
	b, err := json.Marshal(layer)
	if err != nil {
		t.Fatal(err)
	}
	u, _ := url.Parse("https://google.com/hello/0")
	driller, err := esri.GetDriller(*u, b)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := driller.(*esri.Layer); !ok {
		t.Fatal("Driller returned is not a layer")
	}
}
