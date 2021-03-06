package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/apiorno/UrlShortener/data"
	"github.com/apiorno/UrlShortener/model"
	"github.com/gorilla/mux"
)

// URLAssociationsController represents the controller to handle requests for URLs
type URLAssociationsController struct {
	l *log.Logger
}

// Startup initializes the controller
func Startup(l *log.Logger) *mux.Router {
	mux := mux.NewRouter()
	controller := &URLAssociationsController{l}

	getRouter := mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/{id:[0-9a-v]{20}}", controller.RedirectIfExist())
	getRouter.HandleFunc("/", controller.GetAllAssociatedURLs())

	postRouter := mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", controller.AssociateURL())

	deleteRouter := mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9a-v]{20}}", controller.DisassociateURL())

	putRouter := mux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9a-v]{20}}", controller.UpdateURL())

	return mux
}

// LogLine logs the text using the controller's logger
func (c *URLAssociationsController) LogLine(text string) {
	c.l.Println(text)
}

// GetAllAssociatedURLs return all the URLs associations
func (c *URLAssociationsController) GetAllAssociatedURLs() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle GET")

		urlAssocs, err := data.FindAllURLs()
		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Can not retrieve data from firebase", http.StatusInternalServerError)
		}
		err = json.NewEncoder(rw).Encode(urlAssocs)

		if err != nil {
			http.Error(rw, "Can not convert Url Associations to JSON", http.StatusInternalServerError)
		}
	}
}

// RedirectIfExist redirects the request id to its associated real url if exists
func (c *URLAssociationsController) RedirectIfExist() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle GET/{id} ")

		id := mux.Vars(r)["id"]

		urlAssoc, err := data.FindURL(id)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(rw, r, urlAssoc.URL, http.StatusMovedPermanently)
	}
}

//AssociateURL generates a unique id to associate to the URL if it is a valid URL
func (c *URLAssociationsController) AssociateURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle POST ")

		urlAssociation := &model.URLAssociation{}
		err := urlAssociation.FromJSON(r.Body)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}

		_, err = url.ParseRequestURI(urlAssociation.URL)
		if err != nil {
			http.Error(rw, "Invalid url format", http.StatusBadRequest)
			return
		}
		urlAssociation, err = data.AssociateURL(urlAssociation.URL)
		if err != nil {
			http.Error(rw, "Can not associate url", http.StatusInternalServerError)
		}
		err = urlAssociation.ToJSON(rw)

		if err != nil {
			http.Error(rw, "Can not convert to Url Association to JSON", http.StatusInternalServerError)
		}
	}
}

// DisassociateURL removes the association of the URL to the requested id if exists
func (c *URLAssociationsController) DisassociateURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle DELETE/{id} ")

		id := mux.Vars(r)["id"]

		err := data.DisassociateURL(id)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

// UpdateURL updates the url associated to the requested uuid
func (c *URLAssociationsController) UpdateURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.LogLine("Handle UPDATE/{id} ")

		id := mux.Vars(r)["id"]

		urlAssociation := &model.URLAssociation{}
		err := urlAssociation.FromJSON(r.Body)

		if err != nil {
			c.LogLine(err.Error())
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}

		_, err = url.ParseRequestURI(urlAssociation.URL)
		if err != nil {
			http.Error(rw, "Invalid url format", http.StatusBadRequest)
			return
		}

		err = data.UpdateURL(id, urlAssociation.URL)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
