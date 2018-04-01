// Package bapi, functions to connect with Iota IRI to find addresses/transactions/bundles/tags
package bapi

import (
	"fmt"

	"bitbucket.org/thex/iotaBrowser/b/iri"
	"github.com/iotaledger/giota"
)

// GetNoneInfo returns the node information
func IriNodeInfo() (info, milestone string, err error) {
	var ni *giota.GetNodeInfoResponse
	var url string
	url, ni, err = iri.GetNodeInfo()
	if err != nil {
		return info, milestone, err
	}
	return ni.AppName + " " + ni.AppVersion + " " + url, string(ni.LatestSolidSubtangleMilestone), err
}

// Transaction is the response transaction struct used to return a Iota transaction
type Transaction struct {
	Signature       string
	Address         string
	Value           int64
	ObsoleteTag     string
	Timestamp       string
	CurrentIndex    int64
	LastIndex       int64
	Bundle          string
	TrunkTrx        string
	BranchTrx       string
	Tag             string
	AttachmentTs    string
	AttachmentTsLow string
	AttachmentTsUp  string
	Nonce           string
}

// newTransaction sets a Transaction structure with a giota.Transaction
func newTransaction(s *giota.Transaction, tran *Transaction) {
	*tran = Transaction{
		Signature:       string(s.SignatureMessageFragment),
		Address:         string(s.Address),
		Value:           s.Value,
		ObsoleteTag:     string(s.ObsoleteTag),
		Timestamp:       s.Timestamp.String(),
		CurrentIndex:    s.CurrentIndex,
		LastIndex:       s.LastIndex,
		Bundle:          string(s.Bundle),
		TrunkTrx:        string(s.TrunkTransaction),
		BranchTrx:       string(s.BranchTransaction),
		Tag:             string(s.Tag),
		AttachmentTs:    string(s.AttachmentTimestamp),
		AttachmentTsLow: string(s.AttachmentTimestampLowerBound),
		AttachmentTsUp:  string(s.AttachmentTimestampUpperBound),
	}
}

// IriTangleAddr returns an address balance and transactions
func IriTangleAddr(addrStr string, trans *[]Transaction) (bal int64, hashes []string, states []bool, err error) {
	fmt.Println("IriTangleAddr", addrStr)
	addr := giota.Address(addrStr)

	// address balance
	addrs := []giota.Address{addr}
	bals, err := iri.GetBalances(addrs)
	if bals == nil {
		return bal, hashes, states, err
	}
	bal = bals[0]
	// address transactions
	var hashesT []giota.Trytes
	hashesT, states, err = iri.GetTransactions(true, addrs, nil, nil, nil)
	if err == nil && hashesT != nil && len(hashesT) > 0 {
		hashes = make([]string, len(hashesT))
		for k, v := range hashesT {
			hashes[k] = string(v)
		}
		// address transactions detail
		var trxs *[]giota.Transaction
		trxs, _, _, err = iri.GetTrytes(false, false, hashesT)
		if err != nil {
			return bal, hashes, states, err
		}
		// create the transaction array, point the array to trans (so it is returned the data directly)
		sa := make([]Transaction, len(*trxs))
		*trans = sa
		for i := 0; i < len(sa); i++ {
			newTransaction(&(*trxs)[i], &sa[i])
		}
	}
	return bal, hashes, states, err
}

// IriTangleHash shows a transaction detail and children transactions
func IriTangleHash(hashStr string, trans *[]Transaction) (hashes []string, states []bool, err error) {
	fmt.Println("IriTangleHash", hashStr)
	hash := giota.Trytes(hashStr)

	// transaction detail
	hashesT := []giota.Trytes{giota.Trytes(hash)}
	var trxs *[]giota.Transaction
	trxs, states, _, err = iri.GetTrytes(true, false, hashesT)
	if err != nil {
		return hashes, states, err
	}
	sa := make([]Transaction, 1)
	*trans = sa
	newTransaction(&(*trxs)[0], &sa[0])

	// transaction children
	hashesT, _, err = iri.GetTransactions(false, nil, nil, nil, hashesT)
	if err == nil && hashesT != nil && len(hashesT) > 0 {
		hashes = make([]string, len(hashesT))
		for k, v := range hashesT {
			hashes[k] = string(v)
		}
	}
	return hashes, states, err
}

// IriTangleBundle returns the transactions in a bundle
func IriTangleBundle(bundStr string, trans *[]Transaction) (hashes []string, states []bool, err error) {
	fmt.Println("IriTangleBund", bundStr)
	bund := giota.Trytes(bundStr)

	// bundle transactions
	bundsT := []giota.Trytes{bund}
	var hashesT []giota.Trytes
	hashesT, states, err = iri.GetTransactions(true, nil, nil, bundsT, nil)
	if err == nil && hashesT != nil && len(hashesT) > 0 {
		hashes = make([]string, len(hashesT))
		for k, v := range hashesT {
			hashes[k] = string(v)
		}
		// address transactions detail
		var trxs *[]giota.Transaction
		trxs, states, _, err = iri.GetTrytes(true, false, hashesT)
		if err != nil {
			return hashes, states, err
		}
		// create the transaction array, point the array to trans (so it is returned the data directly)
		sa := make([]Transaction, len(*trxs))
		*trans = sa
		for i := 0; i < len(sa); i++ {
			newTransaction(&(*trxs)[i], &sa[i])
		}
	}
	return hashes, states, err
}

// IriTangleObsoleteTag returns the transactions with the obsoleteTag requested
func IriTangleObsoleteTag(obsTagStr string) (hashes []string, err error) {
	fmt.Println("IriTangleOtag", obsTagStr)
	obsTag := giota.Trytes(obsTagStr)

	// tag transactions
	obsTagsT := []giota.Trytes{obsTag}
	var hashesT []giota.Trytes
	hashesT, _, err = iri.GetTransactions(false, nil, obsTagsT, nil, nil)
	if err == nil && hashesT != nil && len(hashesT) > 0 {
		hashes = make([]string, len(hashesT))
		for k, v := range hashesT {
			hashes[k] = string(v)
		}
	}
	return hashes, err
}
