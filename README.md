# iotaBrowser
Iota tangle browser developed in Go using giota package

If you are a go developer trying to do something with Iota, good luck, documentation and support is awful.
For instance, as of today there is no testnet tangle browser working with realtime tesntet data.
I hope this code can help you.

This tangle browser consists of a a go web server that listens for the following requests:
- address: returns the address balance and related transactions
- transaction: returns the transaction detail and children transactions
- bundle: returns the bundle transactions
- obsoleteTag: retuns the transactions that have the obsolete tag
(I use obsolete tag because IRI api only returns obsolete tag... yes Iota has obsolete data, don't be surprised...)

To use, first start the server typing in the terminal: go run cmd/server/main.go
There are several flags that you can change to customize, in cmd/server/main.go you have0:
 - domain := flag.String("domain", "localhost", "Server domain")
 - port   := flag.String("port", "3030", "Server port")
 - static := flag.String("static", "/home/gugu/1devel/go/src/bitbucket.org/thex/iotaBrowser/static", "Static file path")
 - iriUrl := flag.String("iriurl", "http://p101.iotaledger.net:14700", "Iota IRI server:port")

ATTENTION: the most important flag is static , depending on you installed this program you must change the path otherwise the server wont find the static folder.
For example to run on port 80 and with static in "/var/www", type: go run cmd/server/main.go -port=3030 - static="/var/www"

Then type in your browser: http://localhost:3030/static/index
The program has an http server that provides the index page to the browser and then you can browse the tangle.
The browsing is done with ajax but the displaying is very basic but works (I have no time for css, if you want to do it you are welcome).

Enjoy the tangle ;-)