// Package iri, connects with Iota IRI
package iri

import (
	"errors"
	"fmt"

	"github.com/iotaledger/giota"
)

type iotaContext struct {
	err      error
	url      string
	api      *giota.API
	security int
	mwm      int64 // minimum weight magnitude
	seed     giota.Trytes
	solMil   giota.Trytes
}

// ic is the Iota IRI context
var ic iotaContext

// Init initiates the IRI adapter
func Init(iriStr *string) {
	ic = iotaContext{
		url: *iriStr,
		api: giota.NewAPI(*iriStr, nil),
	}
	// Call to get the millestone
	GetNodeInfo()
}

func iriError(err error, errStr, fnStr string) error {
	if err != nil {
		ic.err = err
	}
	if errStr != "" {
		ic.err = errors.New(errStr)
	}
	fmt.Println(fnStr + " " + err.Error())
	return ic.err
}

// GetNoneInfo returns the node information
func GetNodeInfo() (url string, ni *giota.GetNodeInfoResponse, err error) {
	// get node info and latest solid milestone ex.. "HMRJFHJDYDCVDZGXJDUZAOVRIPKWBB9XRQKVOMMIWSGJIWORBNENDL9GINCNFQX9MHTWQVNZDKI9IW999"
	url = ic.url
	ni, err = ic.api.GetNodeInfo()
	if err != nil {
		iriError(err, "", "GetNodeInfo")
		return url, ni, err
	}
	ic.solMil = ni.LatestSolidSubtangleMilestone
	return url, ni, err
}

// GetBalances returns the addresses balances
func GetBalances(addrs []giota.Address) (bal []int64, err error) {
	if addrs == nil || len(addrs) == 0 {
		err = iriError(nil, "No addresses", "GetBalances")
		return bal, err
	}
	bals, err := ic.api.GetBalances(addrs, 100)
	if err != nil {
		iriError(err, "", "GetBalances")
		return bal, err
	}
	return bals.Balances, err
}

// GetTransactions gets transactions by address/tags/bundle/approvees, if getInclusion set it also returns the transaction inclusion states
// NOTE Iota bullshit, caution with the obsolete tag, IRI findTransaction only finds the obsolete tag (but the tangle explorer shows the real tag)
func GetTransactions(getInclusion bool, addrs []giota.Address, tags, bunds, aprs []giota.Trytes) (hashes []giota.Trytes, states []bool, err error) {
	ftr, has := giota.FindTransactionsRequest{}, false
	if addrs != nil && len(addrs) > 0 {
		ftr.Addresses, has = addrs, true
	}
	if bunds != nil && len(bunds) > 0 {
		ftr.Bundles, has = bunds, true
	}
	if tags != nil && len(tags) > 0 {
		ftr.Tags, has = tags, true
	}
	if aprs != nil && len(aprs) > 0 {
		ftr.Approvees, has = aprs, true
	}
	if !has {
		err = iriError(nil, "No addresses, tags, approvees or bundles", "GetTransactions")
		return hashes, states, err
	}
	// find transaction hashes
	ft, err := ic.api.FindTransactions(&ftr)
	if err != nil {
		iriError(err, "", "GetTransactions")
		return hashes, states, err
	}
	hashes = ft.Hashes
	if len(hashes) == 0 {
		return hashes, states, err
	}
	// if inclusion state requested, get transaction inclusion
	if getInclusion {
		// check transactions inclusion state
		gis, err := ic.api.GetInclusionStates(hashes, []giota.Trytes{ic.solMil})
		if err != nil {
			iriError(err, "", "GetTransactions inclusion")
			return hashes, states, err
		}
		states = gis.States
	}
	return hashes, states, err
}

// GetTrytes gets transactions trytes
func GetTrytes(getInclusion bool, bund bool, tryts []giota.Trytes) (trans *[]giota.Transaction, states []bool, bunds []giota.Trytes, err error) {
	if tryts == nil || len(tryts) == 0 {
		err = iriError(nil, "No trytes", "GetTrytes")
		return trans, states, bunds, err
	}
	var gt *giota.GetTrytesResponse
	gt, err = ic.api.GetTrytes(tryts)
	if err != nil {
		iriError(err, "", "GetTrytes")
		return trans, states, bunds, err
	}
	trans = &gt.Trytes
	// if inclusion state requested, get transaction inclusion
	if getInclusion {
		// check transactions inclusion state
		gis, err := ic.api.GetInclusionStates(tryts, []giota.Trytes{ic.solMil})
		if err != nil {
			iriError(err, "", "GetTrytes inclusion")
			return trans, states, bunds, err
		}
		states = gis.States
	}
	// if bundles requested return bundles found
	if bund {
		mb := make(map[string]bool)
		for i := 0; i < len(*trans); i++ {
			mb[string((*trans)[i].Bundle)] = true
		}
		bunds = make([]giota.Trytes, len(mb))
		i := 0
		for k, _ := range mb {
			bunds[i] = giota.Trytes(k)
			i++
		}
	}
	return trans, states, bunds, err
}
