package esri

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	// FormatJSON is a query param for formatting an esri request to json
	FormatJSON = "f=json"
)

type urlSetter struct {
	URL url.URL
}

func (u *urlSetter) SetURL(URL url.URL) {
	u.URL = URL
}

// Root response from esri
type Root struct {
	urlSetter
	CurrentVersion float64   `json:"currentVersion"`
	Folders        []string  `json:"folders"`
	Services       []Service `json:"services"`
}

// GetChildren Returns folders & services
func (r Root) GetChildren() []url.URL {
	children := make([]url.URL, 0)
	for _, f := range r.Folders {
		folder := Folder{Name: f}
		children = append(children, folder.GetURL(r.URL))
	}
	for _, s := range r.Services {
		children = append(children, s.GetURL(r.URL))
	}
	return children
}

// CanDrill into a root ESRI request
func (r Root) CanDrill() bool {
	return true
}

// GetURL returns the parent url, since this is the root request
func (r Root) GetURL(parent url.URL) url.URL {
	return parent
}

// Matches if map contains folders, services & currentVersion
func (r Root) Matches(m map[string]interface{}) bool {
	return containsAll([]string{"folders", "services", "currentVersion"}, m)
}

//Folder represents a folder json response from ESRI
type Folder struct {
	urlSetter
	Name           string
	CurrentVersion float64   `json:"currentVersion"`
	Folders        []string  `json:"folders"`
	Services       []Service `json:"services"`
}

//CanDrill into a folder
func (f Folder) CanDrill() bool {
	return true
}

//GetChildren returns services & other folders under this folder
func (f Folder) GetChildren() []url.URL {
	children := make([]url.URL, 0)
	for _, folder := range f.Folders {
		folder := Folder{Name: folder}
		children = append(children, folder.GetURL(f.URL))
	}
	for _, service := range f.Services {
		children = append(children, service.GetURL(f.URL))
	}
	return children
}

// GetURL appends its Name field to existing url.Path
func (f Folder) GetURL(parent url.URL) url.URL {
	return url.URL{
		Scheme:   parent.Scheme,
		Host:     parent.Host,
		Path:     fmt.Sprintf("%s/%s", parent.Path, f.Name),
		RawQuery: FormatJSON,
	}
}

// Matches if contains services & folders
func (f Folder) Matches(m map[string]interface{}) bool {
	return containsAll([]string{"services", "folders", "currentVersion"}, m)
}

// Extract the last part of URL, since that's the name of the Folder
// i.e. arcgis/rest/services/Environment, where Environment is the folder name
func (f *Folder) Extract(req url.URL) {
	f.URL = req
	pathPieces := strings.Split(req.Path, "/")
	f.Name = pathPieces[len(pathPieces)-1]
}

// Service represents a json response for a service from ESRI
type Service struct {
	urlSetter
	Name                  string        `json:"name"`
	Type                  string        `json:"type"`
	CurrentVersion        float64       `json:"currentVersion"`
	ServiceDescription    string        `json:"serviceDescription"`
	MapName               string        `json:"mapName"`
	Description           string        `json:"description"`
	CopyrightText         string        `json:"copyrightText"`
	SupportsDynamicLayers bool          `json:"supportsDynamicLayers"`
	Layers                []Layer       `json:"layers"`
	Tables                []interface{} `json:"tables"`
	SpatialReference      struct {
		Wkid       int `json:"wkid"`
		LatestWkid int `json:"latestWkid"`
	} `json:"spatialReference"`
	SingleFusedMapCache bool `json:"singleFusedMapCache"`
	InitialExtent       struct {
		Xmin             float64 `json:"xmin"`
		Ymin             float64 `json:"ymin"`
		Xmax             float64 `json:"xmax"`
		Ymax             float64 `json:"ymax"`
		SpatialReference struct {
			Wkid       int `json:"wkid"`
			LatestWkid int `json:"latestWkid"`
		} `json:"spatialReference"`
	} `json:"initialExtent"`
	FullExtent struct {
		Xmin             float64 `json:"xmin"`
		Ymin             float64 `json:"ymin"`
		Xmax             float64 `json:"xmax"`
		Ymax             float64 `json:"ymax"`
		SpatialReference struct {
			Wkid       int `json:"wkid"`
			LatestWkid int `json:"latestWkid"`
		} `json:"spatialReference"`
	} `json:"fullExtent"`
	MinScale                  int    `json:"minScale"`
	MaxScale                  int    `json:"maxScale"`
	Units                     string `json:"units"`
	SupportedImageFormatTypes string `json:"supportedImageFormatTypes"`
	DocumentInfo              struct {
		Title                string `json:"Title"`
		Author               string `json:"Author"`
		Comments             string `json:"Comments"`
		Subject              string `json:"Subject"`
		Category             string `json:"Category"`
		AntialiasingMode     string `json:"AntialiasingMode"`
		TextAntialiasingMode string `json:"TextAntialiasingMode"`
		Keywords             string `json:"Keywords"`
	} `json:"documentInfo"`
	Capabilities          string `json:"capabilities"`
	SupportedQueryFormats string `json:"supportedQueryFormats"`
	ExportTilesAllowed    bool   `json:"exportTilesAllowed"`
	MaxRecordCount        int    `json:"maxRecordCount"`
	MaxImageHeight        int    `json:"maxImageHeight"`
	MaxImageWidth         int    `json:"maxImageWidth"`
	SupportedExtensions   string `json:"supportedExtensions"`
}

// CanDrill into a service
func (s Service) CanDrill() bool {
	return s.Type != "GeometryServer"
}

// GetChildren returns child layers
func (s Service) GetChildren() []url.URL {
	children := make([]url.URL, 0)
	for _, l := range s.Layers {
		children = append(children, l.GetURL(s.URL))
	}
	return children
}

// GetURL returns Path/Name/Type
func (s Service) GetURL(parent url.URL) url.URL {
	name := s.Name
	if strings.Contains(name, "/") {
		splitPath := strings.Split(s.Name, "/")
		name = splitPath[len(splitPath)-1]
	}
	return url.URL{
		Scheme:   parent.Scheme,
		Host:     parent.Host,
		Path:     fmt.Sprintf("%s/%s/%s", parent.Path, name, s.Type),
		RawQuery: FormatJSON,
	}
}

// Matches when map contains proper keys
func (s Service) Matches(m map[string]interface{}) bool {
	return containsAll([]string{"serviceDescription", "layers", "supportsDynamicLayers"}, m)
}

// Extract extracts service type & service name from URL
// Scheme: /path/to/<service.Name>/<service.Type>
func (s *Service) Extract(req url.URL) {
	pathPieces := strings.Split(req.Path, "/")
	s.Name = pathPieces[len(pathPieces)-2]
	s.Type = pathPieces[len(pathPieces)-1]
}

// Layer represents a layer object returned from ESRI
type Layer struct {
	urlSetter
	CurrentVersion float64     `json:"currentVersion"`
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	Type           string      `json:"type"`
	Description    string      `json:"description"`
	GeometryType   interface{} `json:"geometryType"`
	CopyrightText  string      `json:"copyrightText"`
	ParentLayer    interface{} `json:"parentLayer"`
	SubLayers      []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"subLayers"`
	MinScale                               int           `json:"minScale"`
	MaxScale                               int           `json:"maxScale"`
	DefaultVisibility                      bool          `json:"defaultVisibility"`
	Extent                                 Extent        `json:"extent"`
	HasAttachments                         bool          `json:"hasAttachments"`
	HTMLPopupType                          string        `json:"htmlPopupType"`
	DisplayField                           string        `json:"displayField"`
	TypeIDField                            interface{}   `json:"typeIdField"`
	Fields                                 interface{}   `json:"fields"`
	Relationships                          []interface{} `json:"relationships"`
	CanModifyLayer                         bool          `json:"canModifyLayer"`
	CanScaleSymbols                        bool          `json:"canScaleSymbols"`
	HasLabels                              bool          `json:"hasLabels"`
	Capabilities                           string        `json:"capabilities"`
	SupportsStatistics                     bool          `json:"supportsStatistics"`
	SupportsAdvancedQueries                bool          `json:"supportsAdvancedQueries"`
	SupportedQueryFormats                  string        `json:"supportedQueryFormats"`
	OwnershipBasedAccessControlForFeatures struct {
		AllowOthersToQuery bool `json:"allowOthersToQuery"`
	} `json:"ownershipBasedAccessControlForFeatures"`
	UseStandardizedQueries    bool `json:"useStandardizedQueries"`
	AdvancedQueryCapabilities struct {
		UseStandardizedQueries       bool `json:"useStandardizedQueries"`
		SupportsStatistics           bool `json:"supportsStatistics"`
		SupportsOrderBy              bool `json:"supportsOrderBy"`
		SupportsDistinct             bool `json:"supportsDistinct"`
		SupportsPagination           bool `json:"supportsPagination"`
		SupportsTrueCurve            bool `json:"supportsTrueCurve"`
		SupportsReturningQueryExtent bool `json:"supportsReturningQueryExtent"`
		SupportsQueryWithDistance    bool `json:"supportsQueryWithDistance"`
	} `json:"advancedQueryCapabilities"`
}

// Implement Driller interface for layer since we want to signify our end of the tree

// GetChildren for layers has no children, return empty slice
func (l Layer) GetChildren() []url.URL {
	return []url.URL{}
}

// CanDrill returns false because we cannot drill into a layer, it's the lowest layer of an ESRI tree
func (l Layer) CanDrill() bool {
	return false
}

// GetURL appends ID to parent path
func (l Layer) GetURL(parent url.URL) url.URL {
	return url.URL{
		Scheme: parent.Scheme,
		Host:   parent.Host,
		Path:   fmt.Sprintf("%s/%d", parent.Path, l.ID),
	}
}

// Matches if map contains keys specific to layers
func (l Layer) Matches(m map[string]interface{}) bool {
	keys := []string{"description", "id", "type", "extent"}
	return containsAll(keys, m)
}

// Extent represents a layers extent
type Extent struct {
	Xmin             float64 `json:"xmin"`
	Ymin             float64 `json:"ymin"`
	Xmax             float64 `json:"xmax"`
	Ymax             float64 `json:"ymax"`
	SpatialReference struct {
		Wkid       int `json:"wkid"`
		LatestWkid int `json:"latestWkid"`
	} `json:"spatialReference"`
}

// Wkid returns Wkid that represents this extent
func (e Extent) Wkid() int {
	wkid := e.SpatialReference.LatestWkid
	if wkid == 0 {
		wkid = e.SpatialReference.Wkid
	}
	return wkid
}

func containsAll(keys []string, m map[string]interface{}) bool {
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return false
		}
	}
	return true
}
