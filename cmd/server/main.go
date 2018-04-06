// Package main contains the tangle browser server. Execute it with
// go run cmd/server/main.go -domain=localhost -port=3030
package main

import (
	"flag"

	"bitbucket.org/thex/iotaBrowser/a/server"
	"bitbucket.org/thex/iotaBrowser/b/iri"
)

func main() {
	// load flags
	domain := flag.String("domain", "localhost", "Server domain")
	port := flag.String("port", "3030", "Server port")
	static := flag.String("static", "/home/gugu/1devel/go/src/bitbucket.org/thex/iotaBrowser/static", "Static file path") // NOTE adapt static path to your system location
	iriUrl := flag.String("iriurl", "http://iota.nessys.es:14700", "Iota IRI server:port")                                // NOTE if not works try other from this testnet list.. "http://iota.nessys.es:14700", "https://testnet140.tangle.works:443", "http://p101.iotaledger.net:14700", "http://p103.iotaledger.net:14700"

	//  set IRI server
	iri.Init(iriUrl)
	// start server
	server.Init(domain, port, static)
}
