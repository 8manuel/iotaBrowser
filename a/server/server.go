// Package server contains the http server code to handle all the HTTP REST requests.
package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/thex/iotaBrowser/bapi"
	"github.com/julienschmidt/httprouter"
)

// server is an internal variable to keep the server and being able to shutdown later if the admin handler requests it
var server http.Server

// Init initiates the server and adds the routers to the handler
func Init(domain, port, static *string) {
	// instantiate a new mutex router and a server with this mutex
	r := httprouter.New()
	server = http.Server{
		Addr:    *domain + ":" + *port,
		Handler: r,
	}

	// add static file and admin handler to the router
	r.ServeFiles("/static/*filepath", http.Dir(*static))
	// iota browser handler
	r.GET("/get/:obj/:value", TangleGet)

	// start the server
	println("Starting http server on", *domain, *port, ", static files on", *static)
	err := server.ListenAndServe()
	println("Exiting server... ", err.Error())
}

// requestOk valitates that the request is valid
func requestOk(w http.ResponseWriter, r *http.Request) bool {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return false
	}
	return true
}

// responseJson writes content-type, statuscode, payload encoded as json
func responseJson(err error, w http.ResponseWriter, payload interface{}) {
	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	byteArray, _ := json.Marshal(payload)
	//w.Header().Set("Access-Control-Allow-Origin", "*") // uncomment to allow CORS
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", byteArray)
}

// ResPl is the response payload struct used to return the response values, some fields have the omitempty so they will not be included
type ResPl struct {
	Res    int16              `json:"res, omitempty"`
	ErrS   string             `json:"error, omitempty"`
	Serv   string             `json:"server, omitempty"`
	Mile   string             `json:"milestone, omitempty"`
	Addr   string             `json:"address, omitempty"`
	Tran   string             `json:"transaction, omitempty"`
	Bund   string             `json:"bundle, omitempty"`
	Otag   string             `json:"obsoleteTag, omitempty"`
	Abal   int64              `json:"balance, omitempty"`
	Hashes []string           `json:"hashes, omitempty"`
	Incl   []bool             `json:"included, omitempty"`
	Trans  []bapi.Transaction `json:"trans, omitempty"`
}

// TangleGet GET allows to query the tangle by:
//  - address: returns the address balance and related transactions
//  - transaction: returns the transaction detail, parents and children transactions
//  - bundle: returns the bundle transactions
//  - obsoleteTag: retuns the transactions that have the obsolete tag
func TangleGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !requestOk(w, r) {
		return
	}
	var err error
	rpl := ResPl{}
	obj, value := p.ByName("obj"), p.ByName("value")
	// depending on the case execute a different function
	switch obj {
	case "server":
		rpl.Serv, rpl.Mile, err = bapi.IriNodeInfo()
	case "address":
		rpl.Addr = value
		rpl.Abal, rpl.Hashes, rpl.Incl, err = bapi.IriTangleAddr(value, &rpl.Trans)
	case "transaction":
		rpl.Tran = value
		rpl.Hashes, rpl.Incl, err = bapi.IriTangleHash(value, &rpl.Trans)
	case "bundle":
		rpl.Bund = value
		rpl.Hashes, rpl.Incl, err = bapi.IriTangleBundle(value, &rpl.Trans)
	case "obsoleteTag":
		rpl.Otag = value
		rpl.Hashes, err = bapi.IriTangleObsoleteTag(value)
	}
	if err != nil {
		rpl.ErrS = err.Error()
	}
	responseJson(nil, w, rpl)
}
