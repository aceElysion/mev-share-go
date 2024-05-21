package sse

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/aceElysion/mev-share-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Event represents a matchmaker event sent from sse subscription
type Event struct {
	Data  *MatchMakerEvent // Will be nil if an error occurred during poll
	Error error
}

// MatchMakerEvent represents the pending transaction hints sent by matchmaker
type MatchMakerEvent struct {
	Hash        common.Hash          `json:"hash"`
	Logs        []types.Log          `json:"logs,omitempty"`
	Txs         []PendingTransaction `json:"txs,omitempty"`
	MevGasPrice *hexutil.Big         `json:"mevGasPrice,omitempty"`
	GasUsed     *hexutil.Big         `json:"gasUsed,omitempty"`
}

// PendingTransaction represents the hits revealed by the matchmaker about the transaction / bundle
type PendingTransaction struct {
	Hash             *common.Hash    `json:"hash,omitempty"`
	To               *common.Address `json:"to,omitempty"`
	FunctionSelector [4]byte         `json:"functionSelector,omitempty"`
	CallData         []byte          `json:"callData,omitempty"`
}

// UnmarshalJSON unmarshals JSON data into a PendingTransaction
func (t *PendingTransaction) UnmarshalJSON(data []byte) error {
	var temp struct {
		Hash             *common.Hash    `json:"hash,omitempty"`
		To               *common.Address `json:"to"`
		FunctionSelector string          `json:"functionSelector,omitempty"`
		CallData         string          `json:"callData,omitempty"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	t.To = temp.To
	t.Hash = temp.Hash

	if temp.CallData != "" && temp.CallData != "0x" {
		decoded, err := hex.DecodeString(strings.TrimPrefix(temp.CallData, "0x"))
		if err == nil {
			t.CallData = decoded
		}
	}

	if temp.FunctionSelector != "" && temp.FunctionSelector != "0x" {
		decoded, err := hex.DecodeString(strings.TrimPrefix(temp.FunctionSelector, "0x"))
		if err == nil && len(decoded) >= 4 {
			copy(t.FunctionSelector[:], decoded[:4])
		}
	}

	return nil
}
