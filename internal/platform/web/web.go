package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/CRoasSanhez/yofio-test/internal/responses"
	"github.com/jinzhu/gorm"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Server ...
type Server struct {
	router   *httprouter.Router
	DataBase *gorm.DB
	Ctx      context.Context
}

// Service ...
type Service func(s *Server) ServerHandler

// ServerHandler ...
type ServerHandler func(w http.ResponseWriter, r *http.Request, params httprouter.Params) error

// NewServer Return a struct that containts the references to resources
// used in the lifecycle of the API
func NewServer(db *gorm.DB) (*Server, error) {
	srv := Server{
		router:   httprouter.New(),
		DataBase: db,
	}

	return &srv, nil
}

// SetServerCtx ...
func (s *Server) SetServerCtx(ctx context.Context) *Server {
	s.Ctx = ctx
	return s
}

// Handle handle HTTP requests
func (s *Server) Handle(verb, path string, handler ServerHandler) {

	h := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		err := handler(w, r, params)

		if err != nil {
			ResponseJSON(w, http.StatusBadRequest, err)
			return
		}
	}

	s.router.Handle(verb, path, h)
}

// HandleService ...
func (s *Server) HandleService(verb, path string, service Service) {
	h := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		myService := service(s)
		err := myService(w, r, params)
		fmt.Printf("\nService ERR:%v\n", err)

		if err != nil {
			ResponseJSON(w, http.StatusBadRequest, err)
			return
		}
	}

	s.router.Handle(verb, path, h)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Start method starts the internal http server
func (s *Server) Start(appAddr string) {
	err := http.ListenAndServe(appAddr, s.router)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failing starting server")
	}
}

// ResponseJSON returns objects as json
// for httpCode 40x, 50x "response" should be logrus.Fields type
// "respCode" is the struct type RespCode for obtain response message and code
func ResponseJSON(w http.ResponseWriter, httpCode int, response interface{}) {

	if code, ok := response.(responses.RespCode); ok {
		// General Error response
		response = &struct {
			Errors []interface{} `json:"errors"`
		}{
			Errors: []interface{}{code},
		}
		httpCode = code.HTTPStatus
	}

	// Marshal response
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		w.Write(res)
	}
}

// DecodeJSONRequest parse body in to out struct
func DecodeJSONRequest(r *http.Request, out interface{}) error {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &out)
}
