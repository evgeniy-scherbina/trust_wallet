# trust_wallet

## Description:
Project consists of 3 main components:
- custom eth JSON-RPC client
- in-memory storage (which can be replaced with persistent storage in the future)
- parser itself which is backed by previous two components

## Install dependencies (one library for tests: github.com/stretchr/testify)
`go mod tidy`

## Run tests
`go test -v`

## Test description:
- `TestETHClient` - basic test for eth JSON-RPC client
- `TestScanBlockRange` - testing possibility to subscribe to address and then scanning specific block range. Then veryfying that all necessary txs in storage.
- `TestParserInRealTimeMode` - testing parser in real mode (scanning blocks one by one). Then veryfying that all necessary txs in storage.
- `TestParserInRealTimeModeV2` - similar to previous one, but more complex. Storage must contain different txs from different blocks.

NOTE: we using such hack:
```go
inMemStorage := storage.NewInMemoryStorage()
inMemStorage.SetLastProcessedBlock(blockNum - 1)
```
to simplify testing (because we don't know which txs will be in the future)
