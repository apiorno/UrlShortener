package model

import (
	"encoding/json"
	"io"
)

// URLAssociation represents the association of a short url to a real url
type URLAssociation struct {
	UUID string `json:"uuid" firestore:"uuid"`
	URL  string `json:"url" firestore:"url"`
}

// ToJSON generates a json representation of URLAssociation
func (u *URLAssociation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// FromJSON generates an URLAssociation from a json
func (u *URLAssociation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}
