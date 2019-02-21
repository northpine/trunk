package models

import "net/url"

// Driller is some kind of object that can be traversed
type Driller interface {
	// Whether or not we can drill down this server further
	CanDrill() bool
	// Gets urls of resources "children"
	GetChildren() []url.URL
	// Returns a url based on the parents url. This is the url of the request that created this object, not the url of the initial path
	GetURL(parent url.URL) url.URL
	//Sets the URL for the object
	SetURL(URL url.URL)
}

// Matcher matches structs to json blobs. Allows for more strict checking than json.Unmarshal
type Matcher interface {
	Matches(map[string]interface{}) bool
}

// Extracter extracts data from URL to struct. Useful if important data about an object is only found in the URL
type Extracter interface {
	Extract(url.URL)
}

// Sqler turns itself into a SQL insert statement
// Input is starting index of a Psql prepared statement
// Return string is sql statement, returning an interface array of data
type Sqler func(int) (string, []interface{})
