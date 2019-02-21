package esri

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/northpine/trunk/models"
)

type drillerMatcher interface {
	models.Driller
	models.Matcher
}

// GetDriller determines type of ESRI object, and converts that into a driller object
// Will also potentially convert that driller into an extracter to extract useful data from parent object
func GetDriller(URL url.URL, b []byte) (models.Driller, error) {

	m := make(map[string]interface{}, 0)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	matchers := []drillerMatcher{
		&Layer{},
		&Service{},
		&Folder{},
		&Root{},
	}
	for _, matcher := range matchers {
		if matcher.Matches(m) {
			err := json.Unmarshal(b, matcher)
			if err != nil {
				return nil, err
			}
			matcher.SetURL(URL)
			return matcher, nil
		}
	}
	return nil, fmt.Errorf("body failed to unmarshal into a driller: %s", string(b))
}

// AsLayer returns a true, and the layer if the layer is a layer
func AsLayer(driller models.Driller) (*Layer, bool) {
	switch driller.(type) {
	case *Layer:
		l := driller.(*Layer)
		return l, true
	default:
		return nil, false
	}
}

//AsExtracter returns an extracter if the input driller is one
func AsExtracter(driller models.Driller) (models.Extracter, bool) {
	switch driller.(type) {
	case models.Extracter:
		e := driller.(models.Extracter)
		return e, true
	default:
		return nil, false
	}
}
