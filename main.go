package main

import (
	"context"

	lib_eth "github.com/evgeniy-scherbina/trust_wallet/lib/eth"
)

const nodeAddr = "https://cloudflare-eth.com"

type Parser interface {
	// last parsed block
	GetCurrentBlock() (int, error)

	// add address to observer
	Subscribe(address string)

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []*lib_eth.Transaction

	Start(ctx context.Context)
	SetLastProcessedBlock(int)
	ScanBlockRange(blockToStart, blockToEnd int) error
}

func main() {}
