package main

import (
	"fmt"
	"io"
    "encoding/json"
	"bytes"
	"net/http"
)

// curl -d '{"jsonrpc":"2.0","id":"id","method":"get_waku_v2_debug_v1_info", "params":[]}' --header "Content-Type: application/json" http://localhost:8545

type Payload struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []string      `json:"params"` // or []interface{}
}

type Response struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Result  []string      `json:"result"` // or []interface{}
}

type DebugResult struct {
	ListenStr string
}

func main() {
	fmt.Println("JSON RPC request: get_waku_v2_debug_v1_info")

	data := Payload{
		Jsonrpc: "2.0",
		ID: "id",
		Method: "get_waku_v2_debug_v1_info",
		Params: []string{},
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:8545", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}

	defer resp.Body.Close()
	byt, err := io.ReadAll(resp.Body)

    var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}

	var res = dat["result"].(map[string]interface{})

	fmt.Println("listenStr:", res["listenStr"])
}
