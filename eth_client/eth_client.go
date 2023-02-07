package eth_client

import (
	"bytes"
	"encoding/json"
	lib_eth "github.com/evgeniy-scherbina/trust_wallet/lib/eth"
	"io/ioutil"
	"net/http"
	"strconv"
)

type (
	ETHClient struct {
		nodeAddr string
	}

	request struct {
		Jsonrpc string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		Id      int           `json:"id"`
	}

	getBlockNumberResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  string `json:"result"`
		Id      int    `json:"id"`
	}

	getBlockResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  *Block `json:"result"`
		Id      int    `json:"id"`
	}

	Block struct {
		Number       string                 `json:"number"`
		Transactions []*lib_eth.Transaction `json:"transactions"`
	}
)

func New(nodeAddr string) *ETHClient {
	return &ETHClient{
		nodeAddr: nodeAddr,
	}
}

func (eth *ETHClient) GetBlockNumber() (int, error) {
	body, err := eth.call("eth_blockNumber", []interface{}{})
	if err != nil {
		return 0, err
	}

	var resp getBlockNumberResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}
	blockNumber, err := strconv.ParseInt(resp.Result, 0, 64)
	if err != nil {
		return 0, err
	}

	return int(blockNumber), nil
}

func (eth *ETHClient) GetBlock(blockNum int) (*Block, error) {
	blockNumInHex := strconv.FormatInt(int64(blockNum), 16)
	blockNumInHex = "0x" + blockNumInHex
	body, err := eth.call("eth_getBlockByNumber", []interface{}{blockNumInHex, true})
	if err != nil {
		return nil, err
	}

	var resp getBlockResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func (eth *ETHClient) call(method string, params []interface{}) ([]byte, error) {
	req := request{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      83,
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(reqBody)

	resp, err := http.Post(eth.nodeAddr, "application/json", reader)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
