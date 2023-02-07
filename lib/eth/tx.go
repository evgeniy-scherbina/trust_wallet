package eth

type Transaction struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Hash        string `json:"hash"`
	BlockHash   string `json:"blockHash"`
	BlockNumber string `json:"blockNumber"`
}
