package data

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/apiorno/UrlShortener/model"
	"github.com/rs/xid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// ErrURLNotFound is the common error message when an url is not found
var ErrURLNotFound = fmt.Errorf("URL not found")

// DBClient is the firebase client
var DBClient *firestore.Client
var ctx = context.Background()

// Connect tries to connect to firebase bd and return a client
func Connect() *firestore.Client {
	sa := option.WithCredentialsFile("./service-acc-key.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln("Error while connecting")
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("Error while connecting")
		log.Fatalln(err)
	}
	DBClient = client
	return DBClient

}

// FindAllURLs return all the URLAssociations
func FindAllURLs() ([]*model.URLAssociation, error) {

	var assocs []*model.URLAssociation
	urlAssociations := DBClient.Collection("UrlAssociationsData").Doc("UrlAssociationsDoc").Collection("urlAssociations")
	iter := urlAssociations.Documents(ctx)
	defer iter.Stop() // add this line to ensure resources cleaned up
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err

		}
		var assoc *model.URLAssociation
		err = doc.DataTo(&assoc)
		if err != nil {
			return nil, err
		}
		assocs = append(assocs, assoc)
	}
	return assocs, nil

}

// FindURL finds an URLAssociation by id and returns error if not found
func FindURL(ID string) (*model.URLAssociation, error) {

	var assoc *model.URLAssociation

	iter := DBClient.Collection("UrlAssociationsData").Doc("UrlAssociationsDoc").Collection("urlAssociations").Where("uuid", "==", ID).Limit(1).Documents(ctx)
	defer iter.Stop() // add this line to ensure resources cleaned up
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&assoc)
		if err != nil {
			return nil, err
		}
	}
	return assoc, nil
}

// AssociateURL associates an URL to a unique id for short URL
func AssociateURL(url string) (*model.URLAssociation, error) {
	id := xid.New().String()
	urlAssociation := &model.URLAssociation{
		UUID: id,
		URL:  url,
	}
	_, _, err := DBClient.Collection("UrlAssociationsData").Doc("UrlAssociationsDoc").Collection("urlAssociations").Add(ctx, urlAssociation)
	if err != nil {
		return nil, err
	}

	return urlAssociation, nil

}

// DisassociateURL disassociates the URL associated to the requestd id if exists
func DisassociateURL(ID string) error {

	iter := DBClient.Collection("UrlAssociationsData").Doc("UrlAssociationsDoc").Collection("urlAssociations").Where("uuid", "==", ID).Limit(1).Documents(ctx)
	defer iter.Stop() // add this line to ensure resources cleaned up
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return ErrURLNotFound
		}

		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return ErrURLNotFound
		}
	}
	return nil

}

// UpdateURL updates the url associated to the requested uuid
func UpdateURL(ID string, url string) error {

	iter := DBClient.Collection("UrlAssociationsData").Doc("UrlAssociationsDoc").Collection("urlAssociations").Where("uuid", "==", ID).Limit(1).Documents(ctx)
	defer iter.Stop() // add this line to ensure resources cleaned up
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return ErrURLNotFound
		}

		_, err = doc.Ref.Set(ctx, map[string]interface{}{
			"url": url,
		}, firestore.MergeAll)

		if err != nil {
			return ErrURLNotFound
		}
	}
	return nil

}
