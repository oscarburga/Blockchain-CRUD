package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func main_test() {
	testGetVoters()
}

func testGetCandidates() {
	frame := Frame{}
	frame.Cmd = "server"
	frame.Sender = "tester"
	frame.Data = []string{"GET CANDIDATES"}
	send("localhost:8001", frame, func(cn net.Conn) {
		dec := json.NewDecoder(cn)
		var cands []CANDIDATE
		dec.Decode(&cands)
		data, _ := json.MarshalIndent(cands, "", "\t")
		fmt.Println(string(data))
	})

}

func testGetVoters() {
	frame := Frame{}
	frame.Cmd = "server"
	frame.Sender = "tester"
	frame.Data = []string{"GET VOTANTES"}
	send("localhost:8001", frame, func(cn net.Conn) {
		dec := json.NewDecoder(cn)
		var cands []VOTANTE
		dec.Decode(&cands)
		data, _ := json.MarshalIndent(cands, "", "\t")
		fmt.Println(string(data))
	})

}
