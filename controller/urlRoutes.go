package controller

import (
	"net/http"
	"log"
	"github.com/apiorno/UrlShortener/data"
	"github.com/gorilla/mux"
)

// URLAssociationsController represents the controller to handle requests for URLs
type URLAssociationsController struct {
	l *log.Logger
}

//Startup initializes the controller
func Startup(l *log.Logger) *mux.Router {
	mux := mux.NewRouter()
	controller := &URLAssociationsController{l}
	getRouter := mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/{id:[0-9a-v]{20}}", controller.redirectIfExist())
	
	return mux
}

func (c *URLAssociationsController) redirectIfExist() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		urlAssoc, err := data.FindURL(id)

		if err != nil {
			http.Error(rw,err.Error(),http.StatusBadRequest)
			return
		}

		http.Redirect(rw,r,urlAssoc.URL,http.StatusMovedPermanently)
	}
}