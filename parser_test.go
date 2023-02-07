package main

import (
	"context"
	"github.com/evgeniy-scherbina/trust_wallet/parser/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/evgeniy-scherbina/trust_wallet/eth_client"
	"github.com/evgeniy-scherbina/trust_wallet/parser"
)

const (
	addrToTest = "0x8606de07aa34505ff5c2348c60f8cde608340f7c"

	// this block contains tx which has addrToTest as tx.From
	blockNum = 16578442
)

func TestETHClient(t *testing.T) {
	ethClient := eth_client.New(nodeAddr)
	blockNumber, err := ethClient.GetBlockNumber()
	require.NoError(t, err)
	require.Greater(t, blockNumber, 16578493)
}

func TestScanBlockRange(t *testing.T) {
	// setup parser
	ethClient := eth_client.New(nodeAddr)
	var parser Parser = parser.New(ethClient, storage.NewInMemoryStorage())

	// subscribe to address
	addrToTest := "0x8606de07aa34505ff5c2348c60f8cde608340f7c"
	parser.Subscribe(addrToTest)
	// scan block range with only one block
	blockNum := 16578442
	err := parser.ScanBlockRange(blockNum, blockNum)
	require.NoError(t, err)

	txs := parser.GetTransactions(addrToTest)
	require.Len(t, txs, 1)
	tx := txs[0]

	require.Equal(t, addrToTest, tx.From)
	require.Equal(t, "0x3711f702f6550290c07cbe2ca2fb874df6005937", tx.To)
	require.Equal(t, "0x0", tx.Value)
	require.Equal(t, "0x2", tx.Type)
	require.Equal(t, "0x21e63dc4bb94e9c5d67baa8a7ae33f4d770129aaead830bf22df96e38224e579", tx.Hash)
	require.Equal(t, "0x5c17c98430a5f5e526f618c770190786977e56d44f943e130548d2ca17b10da2", tx.BlockHash)
	require.Equal(t, "0xfcf78a", tx.BlockNumber)
}

func TestParserInRealTimeMode(t *testing.T) {
	// setup parser
	ethClient := eth_client.New(nodeAddr)
	var parser Parser = parser.New(ethClient, storage.NewInMemoryStorage())
	parser.SetLastProcessedBlock(blockNum - 1)

	// subscribe to address
	addrToTest := "0x8606de07aa34505ff5c2348c60f8cde608340f7c"
	parser.Subscribe(addrToTest)
	// start parser in real time mode
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go parser.Start(ctx)

	time.Sleep(time.Second * 20)

	txs := parser.GetTransactions(addrToTest)
	require.GreaterOrEqual(t, len(txs), 1)
	tx := txs[0]

	require.Equal(t, addrToTest, tx.From)
	require.Equal(t, "0x3711f702f6550290c07cbe2ca2fb874df6005937", tx.To)
	require.Equal(t, "0x0", tx.Value)
	require.Equal(t, "0x2", tx.Type)
	require.Equal(t, "0x21e63dc4bb94e9c5d67baa8a7ae33f4d770129aaead830bf22df96e38224e579", tx.Hash)
	require.Equal(t, "0x5c17c98430a5f5e526f618c770190786977e56d44f943e130548d2ca17b10da2", tx.BlockHash)
	require.Equal(t, "0xfcf78a", tx.BlockNumber)
}
