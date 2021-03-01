package data

import (
	"fmt"

	"github.com/apiorno/UrlShortener/model"
)

// ErrURLNotFound is the common error message when an url is not found
var ErrURLNotFound = fmt.Errorf("URL not found")

// FindURL finds an URLAssociation by id and returns error if not found
func FindURL(ID string) (*model.URLAssociation, error) {
	for _, u := range urlsAssocList {
		if u.UUID == ID {
			return u, nil
		}
	}
	return nil, ErrURLNotFound
}

var urlsAssocList = []*model.URLAssociation{
	{
		UUID: "abcdefghi12345678910",
		URL:  "http://google.com",
	},
	{
		UUID: "abcdefghi12345678911",
		URL:  "https://facebook.com",
	},
}
