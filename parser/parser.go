package parser

import (
	"context"
	"log"
	"time"

	"github.com/evgeniy-scherbina/trust_wallet/eth_client"
	lib_eth "github.com/evgeniy-scherbina/trust_wallet/lib/eth"
)

type Storage interface {
	LastProcessedBlock() int
	SetLastProcessedBlock(lastProcessedBlock int)
	Subscribe(address string)
	GetTransactions(address string) []*lib_eth.Transaction
	ProcessTx(tx *lib_eth.Transaction)
}

type Parser struct {
	ethClient *eth_client.ETHClient

	storage Storage
}

func New(ethClient *eth_client.ETHClient, storage Storage) *Parser {
	return &Parser{
		ethClient: ethClient,

		storage: storage,
	}
}

func (p *Parser) Subscribe(address string) {
	p.storage.Subscribe(address)
}

func (p *Parser) GetTransactions(address string) []*lib_eth.Transaction {
	return p.storage.GetTransactions(address)
}

func (p *Parser) GetCurrentBlock() (int, error) {
	return p.ethClient.GetBlockNumber()
}

func (p *Parser) Start(ctx context.Context, tickerInternal time.Duration) {
	ticker := time.NewTicker(tickerInternal)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentBlockNumber, err := p.ethClient.GetBlockNumber()
			if err != nil {
				log.Print(err)
				continue
			}

			err = p.ScanBlockRange(p.storage.LastProcessedBlock()+1, currentBlockNumber)
			if err != nil {
				log.Print(err)
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

func (p *Parser) ScanBlockRange(blockToStart, blockToEnd int) error {
	if blockToStart > blockToEnd {
		// nothing to scan
		return nil
	}

	for blockNum := blockToStart; blockNum <= blockToEnd; blockNum++ {
		if err := p.scanBlock(blockNum); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) scanBlock(blockNum int) error {
	block, err := p.ethClient.GetBlock(blockNum)
	if err != nil {
		return err
	}

	for _, tx := range block.Transactions {
		p.storage.ProcessTx(tx)
	}

	p.storage.SetLastProcessedBlock(blockNum)

	//fmt.Printf("Processed Transactions: %v\n", len(block.Transactions))

	return nil
}
