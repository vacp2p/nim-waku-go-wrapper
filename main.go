package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
    "encoding/json"
	"time"
	"bytes"
	"os"
	"os/exec"
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
	cmd := exec.Command("./wakunode2")

    outfile, err := os.Create("./wakunode2.log")
	if err != nil {
		panic(err)
    }
    defer outfile.Close()

	//cmd.Stdout = os.Stdout
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

    fmt.Printf("wakunode2 start, [PID] %d running...\n", cmd.Process.Pid)
    ioutil.WriteFile("wakunode2.lock", []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0666)
	log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)

	time.Sleep(2000 * time.Millisecond)

	// Run this in background
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

	// Stop process
	// Since we have reference to same process we can also use cmd.Process.Kill()
	strb, _ := ioutil.ReadFile("wakunode2.lock")
	command := exec.Command("kill", string(strb))
	command.Start()
	fmt.Println("Stopping wakunode2 process")
}
