# iotaBrowser
Iota tangle browser developed in Go using giota package

If you are a go developer trying to do something with Iota, good luck, documentation and support is awful.
For instance, as of today there is no testnet tangle browser working with realtime tesntet data.
I hope this code can help you.

This tangle browser consists of a a go web server that listens for the following requests:
- address: returns the address balance and related transactions
- transaction: returns the transaction detail, parents and children transactions
- bundle: returns the bundle transactions
- obsoleteTag: retuns the transactions that have the obsolete tag
(I use obsolete tag because IRI api only returns obsolete tag... yes Iota has obsolete data, don't be surprised...)

So far the server will return a json; when I have time I'll try to do some html/javascript to show this properly in the screen.
