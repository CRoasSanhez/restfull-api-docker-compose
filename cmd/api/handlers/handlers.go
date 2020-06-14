package handlers

import (
	"net/http"

	"github.com/CRoasSanhez/yofio-test/internal/platform/web"
	"github.com/julienschmidt/httprouter"
)

func Register(server *web.Server) web.ServerHandler {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
		return nil
	}
}

func Login(server *web.Server) web.ServerHandler {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
		return nil

	}
}

func MembershipPayment(server *web.Server) web.ServerHandler {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
		return nil
	}
}

func MembershipConsult(server *web.Server) web.ServerHandler {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
		return nil

	}
}
