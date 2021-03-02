package view

import (
	"encoding/json"
	"io"
)

//URLPostRequest is an http post request
type URLPostRequest struct {
	URL string `json:"url"`
}

// FromJSON generates an URLAssociation from a json
func (u*URLPostRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}
