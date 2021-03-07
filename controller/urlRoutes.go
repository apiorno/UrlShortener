package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/apiorno/UrlShortener/data"
	"github.com/apiorno/UrlShortener/model"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func incrementRequestOfDuration(d float64) {
	go func() {
		requestsDurationSecondsCount.Inc()
		requestsDurationSecondsSum.Add(d)
		if d <= 10 {
			requestsDurationSecondsBucket.WithLabelValues("10").Inc()
		}
		if d <= 5 {
			requestsDurationSecondsBucket.WithLabelValues("5").Inc()
		}
		if d <= 1 {
			requestsDurationSecondsBucket.WithLabelValues("1").Inc()
		}
		if d <= 0.3 {
			requestsDurationSecondsBucket.WithLabelValues("0.3").Inc()
		}
	}()
}

var (
	requestsDurationSecondsSum = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_request_duration_seconds_sum",
		Help: "Sum of seconds spent on all requests",
	})
)

var (
	requestsDurationSecondsCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_request_duration_seconds_count",
		Help: "Count of  all requests",
	})
)

var (
	requestsDurationSecondsBucket = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_duration_seconds_bucket",
		Help: "group request by tiem repsonses tags",
	},
		[]string{"le"})
)

// URLAssociationsController represents the controller to handle requests for URLs
type URLAssociationsController struct {
	l *log.Logger
}

// Startup initializes the controller
func Startup(l *log.Logger) *mux.Router {
	mux := mux.NewRouter()
	controller := &URLAssociationsController{l}

	requestsDurationSecondsBucket.WithLabelValues("10")
	requestsDurationSecondsBucket.WithLabelValues("5")
	requestsDurationSecondsBucket.WithLabelValues("1")
	requestsDurationSecondsBucket.WithLabelValues("0.3")

	getRouter := mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/{id:[0-9a-v]{20}}", controller.RedirectIfExist())
	getRouter.HandleFunc("/", controller.GetAllAssociatedURLs())

	postRouter := mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", controller.AssociateURL())

	deleteRouter := mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9a-v]{20}}", controller.DisassociateURL())

	putRouter := mux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9a-v]{20}}", controller.UpdateURL())

	mux.Handle("/metrics", promhttp.Handler())
	return mux
}

// LogLine logs the text using the controller's logger
func (c *URLAssociationsController) LogLine(text string) {
	c.l.Println(text)
}

// GetAllAssociatedURLs return all the URLs associations
func (c *URLAssociationsController) GetAllAssociatedURLs() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		c.LogLine("Handle GET")

		urlAssocs, err := data.FindAllURLs()
		if err != nil {
			c.LogLine(err.Error())
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, "Can not retrieve data from firebase", http.StatusInternalServerError)
		}
		err = json.NewEncoder(rw).Encode(urlAssocs)

		if err != nil {
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, "Can not convert Url Associations to JSON", http.StatusInternalServerError)
		}
		incrementRequestOfDuration(float64(time.Now().Sub(t1)))
	}
}

// RedirectIfExist redirects the request id to its associated real url if exists
func (c *URLAssociationsController) RedirectIfExist() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		c.LogLine("Handle GET/{id} ")

		id := mux.Vars(r)["id"]

		urlAssoc, err := data.FindURL(id)

		if err != nil {
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		incrementRequestOfDuration(float64(time.Now().Sub(t1)))
		http.Redirect(rw, r, urlAssoc.URL, http.StatusMovedPermanently)
	}
}

//AssociateURL generates a unique id to associate to the URL if it is a valid URL
func (c *URLAssociationsController) AssociateURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		c.LogLine("Handle POST ")

		urlAssociation := &model.URLAssociation{}
		err := urlAssociation.FromJSON(r.Body)

		if err != nil {
			c.LogLine(err.Error())
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}

		_, err = url.ParseRequestURI(urlAssociation.URL)
		if err != nil {
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, "Invalid url format", http.StatusBadRequest)
			return
		}
		urlAssociation, err = data.AssociateURL(urlAssociation.URL)
		if err != nil {
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, "Can not associate url", http.StatusInternalServerError)
		}
		err = urlAssociation.ToJSON(rw)

		if err != nil {
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, "Can not convert to Url Association to JSON", http.StatusInternalServerError)
		}
		incrementRequestOfDuration(float64(time.Now().Sub(t1)))
	}
}

// DisassociateURL removes the association of the URL to the requested id if exists
func (c *URLAssociationsController) DisassociateURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		c.LogLine("Handle DELETE/{id} ")

		id := mux.Vars(r)["id"]

		err := data.DisassociateURL(id)
		if err != nil {
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		incrementRequestOfDuration(float64(time.Now().Sub(t1)))
	}
}

// UpdateURL updates the url associated to the requested uuid
func (c *URLAssociationsController) UpdateURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		c.LogLine("Handle UPDATE/{id} ")

		id := mux.Vars(r)["id"]

		urlAssociation := &model.URLAssociation{}
		err := urlAssociation.FromJSON(r.Body)

		if err != nil {
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}

		_, err = url.ParseRequestURI(urlAssociation.URL)
		if err != nil {
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, "Invalid url format", http.StatusBadRequest)
			return
		}

		err = data.UpdateURL(id, urlAssociation.URL)

		if err != nil {
			incrementRequestOfDuration(float64(time.Now().Sub(t1)))
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		incrementRequestOfDuration(float64(time.Now().Sub(t1)))
	}
}
