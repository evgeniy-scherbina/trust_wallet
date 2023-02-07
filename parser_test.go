package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/evgeniy-scherbina/trust_wallet/eth_client"
	"github.com/evgeniy-scherbina/trust_wallet/parser"
	"github.com/evgeniy-scherbina/trust_wallet/parser/storage"
)

const (
	addrToTest  = "0x8606de07aa34505ff5c2348c60f8cde608340f7c"
	addrToTest2 = "0x11b01c422b46c139927a697d75c7f466d2f952b8"
	addrToTest3 = "0xeaa900b93c1df26004286362a338370d5b0498f6"

	addrToTestNxtBlock  = "0x1264f83b093abbf840ea80a361988d19c7f5a686"
	addrToTestNxtBlock2 = "0x51922589aa83c3144b054d28d54404c49f0887df"
	addrToTestNxtBlock3 = "0xf18bc040a70f3f192660f5272ddcb5c000b0d7ff"

	// this block contains tx which has addrToTest as tx.From
	blockNum = 16578442

	defaultTickerInternal = time.Second * 5
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
	inMemStorage := storage.NewInMemoryStorage()
	inMemStorage.SetLastProcessedBlock(blockNum - 1)
	var parser Parser = parser.New(ethClient, inMemStorage)

	// subscribe to address
	parser.Subscribe(addrToTest)
	// start parser in real time mode
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go parser.Start(ctx, defaultTickerInternal)

	// give some amount of time to make sure that `blockNum` will be processed
	time.Sleep(2 * defaultTickerInternal)

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

func TestParserInRealTimeModeV2(t *testing.T) {
	// setup parser
	ethClient := eth_client.New(nodeAddr)
	inMemStorage := storage.NewInMemoryStorage()
	inMemStorage.SetLastProcessedBlock(blockNum - 1)
	var parser Parser = parser.New(ethClient, inMemStorage)

	// subscribe to few addresses
	parser.Subscribe(addrToTest)
	parser.Subscribe(addrToTest2)
	parser.Subscribe(addrToTest3)
	parser.Subscribe(addrToTestNxtBlock)
	parser.Subscribe(addrToTestNxtBlock2)
	parser.Subscribe(addrToTestNxtBlock3)
	// start parser in real time mode
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go parser.Start(ctx, defaultTickerInternal)

	// give some amount of time to make sure that `blockNum` and `blockNum+1` will be processed
	time.Sleep(2 * defaultTickerInternal)

	for _, addrToTest := range []string{addrToTest, addrToTest2, addrToTest3, addrToTestNxtBlock, addrToTestNxtBlock2, addrToTestNxtBlock3} {
		txs := parser.GetTransactions(addrToTest)
		require.GreaterOrEqual(t, len(txs), 1)
		tx := txs[0]

		require.Equal(t, addrToTest, tx.From)
	}
}
