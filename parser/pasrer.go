package parser

import (
	"context"
	"fmt"
	"github.com/evgeniy-scherbina/trust_wallet/eth_client"
	lib_eth "github.com/evgeniy-scherbina/trust_wallet/lib/eth"
	"github.com/evgeniy-scherbina/trust_wallet/parser/storage"
	"log"
	"time"
)

type Parser struct {
	ethClient          *eth_client.ETHClient
	lastProcessedBlock int

	storage *storage.InMemoryStorage
}

func New(ethClient *eth_client.ETHClient) *Parser {
	return &Parser{
		ethClient:          ethClient,
		lastProcessedBlock: 0,

		storage: storage.NewInMemoryStorage(),
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

func (p *Parser) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			currentBlockNumber, err := p.ethClient.GetBlockNumber()
			if err != nil {
				log.Print(err)
				continue
			}

			err = p.ScanBlockRange(p.lastProcessedBlock+1, currentBlockNumber)
			if err != nil {
				log.Print(err)
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

func (p *Parser) SetLastProcessedBlock(lastProcessedBlock int) {
	p.lastProcessedBlock = lastProcessedBlock
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

	fmt.Printf("Processed Transactions: %v\n", len(block.Transactions))

	return nil
}
