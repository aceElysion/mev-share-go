package sse

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestPendingTransaction_UnmarshalJSON(t *testing.T) {
	rawJSON := `{
		"to": "0x1234567890abcdef1234567890abcdef12345678",
		"functionSelector": "0xabcdef12",
		"callData": "0xdeadbeef",
		"hash": "0x42bfaa1a6ec6f136a970f102871de727dc3f9c29c0dcc9409fe54de2c60358b1"
	}`

	expectedTo := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	expectedHash := common.HexToHash("0x42bfaa1a6ec6f136a970f102871de727dc3f9c29c0dcc9409fe54de2c60358b1")
	expectedFunctionSelector := [4]byte{0xab, 0xcd, 0xef, 0x12}
	expectedCallData := []byte{0xde, 0xad, 0xbe, 0xef}

	var tx PendingTransaction
	err := json.Unmarshal([]byte(rawJSON), &tx)
	assert.NoError(t, err)
	assert.Equal(t, expectedTo.Cmp(*tx.To), 0)
	assert.Equal(t, expectedHash.Cmp(*tx.Hash), 0)
	assert.Equal(t, expectedFunctionSelector, tx.FunctionSelector)
	assert.Equal(t, expectedCallData, tx.CallData)
}

func TestPendingTransaction_UnmarshalJSON_EmptyFields(t *testing.T) {
	rawJSON := `{
		"to": "0x1234567890abcdef1234567890abcdef12345678",
		"functionSelector": "",
		"callData": "",
		"hash": "0x42bfaa1a6ec6f136a970f102871de727dc3f9c29c0dcc9409fe54de2c60358b1"
	}`

	expectedTo := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	expectedHash := common.HexToHash("0x42bfaa1a6ec6f136a970f102871de727dc3f9c29c0dcc9409fe54de2c60358b1")
	expectedFunctionSelector := [4]byte{}
	var expectedCallData []byte

	var tx PendingTransaction
	err := json.Unmarshal([]byte(rawJSON), &tx)
	assert.NoError(t, err)
	assert.Equal(t, expectedTo.Cmp(*tx.To), 0)
	assert.Equal(t, expectedHash.Cmp(*tx.Hash), 0)
	assert.Equal(t, expectedFunctionSelector, tx.FunctionSelector)
	assert.Equal(t, expectedCallData, tx.CallData)
}

func TestPendingTransaction_UnmarshalJSON_MissingFields(t *testing.T) {
	rawJSON := `{
		"to": "0x1234567890abcdef1234567890abcdef12345678",
		"hash": "0x42bfaa1a6ec6f136a970f102871de727dc3f9c29c0dcc9409fe54de2c60358b1"
	}`

	expectedTo := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	expectedHash := common.HexToHash("0x42bfaa1a6ec6f136a970f102871de727dc3f9c29c0dcc9409fe54de2c60358b1")
	var expectedFunctionSelector [4]byte
	var expectedCallData []byte

	var tx PendingTransaction
	err := json.Unmarshal([]byte(rawJSON), &tx)
	assert.NoError(t, err)
	assert.Equal(t, expectedTo.Cmp(*tx.To), 0)
	assert.Equal(t, expectedHash.Cmp(*tx.Hash), 0)
	assert.Equal(t, expectedFunctionSelector, tx.FunctionSelector)
	assert.Equal(t, expectedCallData, tx.CallData)
}

func TestEvent_Data_Error(t *testing.T) {
	errMsg := "some error"
	event := Event{Error: errors.New(errMsg)}

	assert.Error(t, event.Error)
	assert.Nil(t, event.Data)
}

func TestEvent_Data_Success(t *testing.T) {
	matchMakerEvent := &MatchMakerEvent{
		Hash: common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
	}

	event := Event{Data: matchMakerEvent}

	assert.NoError(t, event.Error)
	assert.Equal(t, matchMakerEvent, event.Data)
}
