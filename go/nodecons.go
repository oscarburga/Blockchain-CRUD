package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

type Frame struct {
	Cmd    string   `json:"cmd"`
	Sender string   `json:"sender"`
	Data   []string `json:"data"`
}

var (
	host         string
	participants int
	chRemotes    chan []string
	chCons       chan ConsensusData
	blockchain   chan *Blockchain
)

func main_nodecons() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		log.Println("Hostname not given")
	} else {
		host = os.Args[1]
		chRemotes = make(chan []string, 1)
		chCons = make(chan ConsensusData)
		blockchain = make(chan *Blockchain, 1)
		if len(os.Args) == 2 {
			createBlockchain()
		}
		chRemotes <- []string{}
		if len(os.Args) >= 3 {
			connectToNode(os.Args[2])
			requestFullBlockChain(os.Args[2])
		}
		server()
	}
}

func createBlockchain() {
	b := CreateTestBlockchain()
	blockchain <- b
}

func connectToNode(remote string) {
	remotes := <-chRemotes
	remotes = append(remotes, remote)
	chRemotes <- remotes
	if !send(remote, Frame{"hello", host, []string{}}, func(cn net.Conn) {
		dec := json.NewDecoder(cn)
		var frame Frame
		dec.Decode(&frame)
		remotes := <-chRemotes
		remotes = append(remotes, frame.Data...)
		chRemotes <- remotes
		log.Printf("%s: friends: %s\n", host, remotes)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}
}

func send(remote string, frame Frame, callback func(net.Conn)) bool {
	if cn, err := net.Dial("tcp", remote); err == nil {
		defer cn.Close()
		enc := json.NewEncoder(cn)
		enc.Encode(frame)
		if callback != nil {
			callback(cn)
		}
		return true
	} else {
		log.Printf("%s: can't connect to %s\n", host, remote)
		idx := -1
		remotes := <-chRemotes
		for i, rem := range remotes {
			if remote == rem {
				idx = i
				break
			}
		}
		if idx >= 0 {
			remotes[idx] = remotes[len(remotes)-1]
			remotes = remotes[:len(remotes)-1]
		}
		chRemotes <- remotes
		return false
	}
}

func server() {
	if ln, err := net.Listen("tcp", host); err == nil {
		defer ln.Close()
		log.Printf("Listening on %s\n", host)
		for {
			if cn, err := ln.Accept(); err == nil {
				go fauxDispatcher(cn)
			} else {
				log.Printf("%s: cant accept connection.\n", host)
			}
		}
	} else {
		log.Printf("Can't listen on %s\n", host)
	}
}

func requestFullBlockChain(remote string) {
	send(remote, Frame{"blockchain", host, []string{}}, func(cn net.Conn) {
		dec := json.NewDecoder(cn)
		var frame Frame
		dec.Decode(&frame)
		copied := MakeBlockchainFromJson([]byte(frame.Data[0]))

		blockchain <- &copied
		fmt.Println("Recieved blockchain:\n", string(copied.GetBlockchainJson()))
	})
}

func fauxDispatcher(cn net.Conn) {
	defer cn.Close()
	dec := json.NewDecoder(cn)
	frame := &Frame{}
	dec.Decode(frame)
	fmt.Println("faux dispatcher with cmd ", frame.Cmd)
	switch frame.Cmd {
	case "hello":
		handleHello(cn, frame)
	case "add":
		handleAdd(frame)
	case "vote":
		ReceiveConsensusVote(frame)
	case "blockchain":
		handleBlockchain(cn)
	case "transaction":
		handleTransaction(cn, frame)
	case "server":
		handleServer(cn, frame)
	}
}

func handleHello(cn net.Conn, frame *Frame) {
	log.Printf("received hello from %s\n", host)
	enc := json.NewEncoder(cn)
	remotes := <-chRemotes
	enc.Encode(Frame{"<response>", host, remotes})
	notification := Frame{"add", host, []string{frame.Sender}}
	for _, remote := range remotes {
		log.Printf("sending add to %s\n", remote)
		send(remote, notification, nil)
	}
	remotes = append(remotes, frame.Sender)
	log.Printf("%s: friends: %s\n", host, remotes)
	chRemotes <- remotes
}

func handleServer(cn net.Conn, frame *Frame) {
	fmt.Println("HANDLE SERVER")
	enc := json.NewEncoder(cn)
	// frame sin data: esta solicitando remotes
	if len(frame.Data) == 0 {
		remotes := <-chRemotes
		enc.Encode(Frame{"<response>", host, remotes})
		log.Printf("Server connection")
		log.Printf("%s: friends: %s\n", host, remotes)
		chRemotes <- remotes
		return
	}

	sendVoters := func(b *Blockchain) {
		b.Database.EnsureInitialized()
		voters := b.Database.Voters
		fmt.Println("Sending voters: ", len(voters))
		enc.Encode(voters)
	}

	sendCandidates := func(b *Blockchain) {
		b.Database.EnsureInitialized()
		cands := b.Database.Candidates
		fmt.Println("Sending candidates: ", len(cands))
		enc.Encode(cands)
	}

	if len(frame.Data) == 1 {
		b := <-blockchain
		if frame.Data[0] == TransCmdGetCandidates {
			sendCandidates(b)
		} else if frame.Data[0] == TransCmdGetVotantes {
			sendVoters(b)
		}
		blockchain <- b
		return
	}

	if len(frame.Data) == 2 {
		trans := Transaction{}
		jsonData := []byte(frame.Data[1])
		trans.Cmd = frame.Data[0]
		trans.Sender = frame.Sender
		bSendCandidates := trans.Cmd == TransCmdCreateCandidate ||
			trans.Cmd == TransCmdUpdateCandidate
		bSendVoters := trans.Cmd == TransCmdCreateVotante ||
			trans.Cmd == TransCmdUpdateVotante
		if bSendCandidates {
			json.Unmarshal(jsonData, &trans.CandidateData)
		} else if bSendVoters {
			json.Unmarshal(jsonData, &trans.VotanteData)
		} else {
			panic("invalid cmd")
		}

		replyToApiCallback := func(wonConsensus bool) {
			// reply to api on consensus finished
			chain := <-blockchain
			fmt.Println("transaction sends: ", bSendCandidates, bSendVoters)
			if bSendCandidates {
				sendCandidates(chain)
			} else if bSendVoters {
				sendVoters(chain)
			}
			blockchain <- chain
		}

		b := <-blockchain
		bDidTransaction := b.AddTransaction(trans)
		if bDidTransaction {
			// send notice for everyone to add a transaction
			msg := Frame{}
			msg.Cmd = "transaction"
			msg.Sender = host
			transData, _ := json.MarshalIndent(trans, "", "\t")
			msg.Data = []string{string(transData)}
			remotes := <-chRemotes
			chRemotes <- remotes
			for _, remote := range remotes {
				send(remote, msg, nil)
			}
			StartConsensus(b.TailBlock.Hash, replyToApiCallback)
			blockchain <- b
		} else {
			blockchain <- b
			replyToApiCallback(false)
		}
	}
}

func handleTransaction(cn net.Conn, frame *Frame) {
	transData := []byte(frame.Data[0])
	trans := Transaction{}
	json.Unmarshal(transData, &trans)
	fmt.Println("received transaction")
	b := <-blockchain
	if b.AddTransaction(trans) {
		StartConsensus(b.TailBlock.Hash, nil)
	}
	blockchain <- b
}

func handleBlockchain(cn net.Conn) {
	enc := json.NewEncoder(cn)
	b := <-blockchain
	jsonData := string(b.GetBlockchainJson())
	frame := Frame{"here you go", host, []string{jsonData}}
	blockchain <- b
	enc.Encode(frame)
}

func handleAdd(frame *Frame) {
	log.Printf("%s: received add from: %s\n", host, frame.Sender)
	remotes := <-chRemotes
	remotes = append(remotes, frame.Data...)
	log.Printf("%s: friends: %s\n", host, remotes)
	chRemotes <- remotes
}

func StartConsensus(blockHash string, onConsensusDone func(bool)) {

	remotes := <-chRemotes
	chRemotes <- remotes
	participants = len(remotes) + 1
	log.Printf("%s: start consensus with remotes %s", host, remotes)

	if participants > 1 {
		consData := ConsensusData{}
		consData.onConsensusDone = onConsensusDone
		consData.AddVote(blockHash)
		consData.myHash = blockHash
		for _, remote := range remotes {
			log.Printf("%s: consensus sending %s to %s\n", host, blockHash, remote)
			send(remote, Frame{"vote", host, []string{blockHash}}, nil)
		}
		remotes = <-chRemotes
		chRemotes <- remotes
		participants = len(remotes) + 1
		chCons <- consData
	}
}

func ReceiveConsensusVote(frame *Frame) {
	vote := frame.Data[0]
	consData := <-chCons
	consData.AddVote(vote)
	log.Printf("%s:received vote from %s - voted %s\n", host, frame.Sender, vote)

	if consData.totalCount < participants {
		log.Printf("%s: participants %d - voted %d\n", host, participants, consData.totalCount)
		chCons <- consData
	} else {
		log.Printf("%s: participants %d - voted %d\n", host, participants, consData.totalCount)
		if consData.maxHash != consData.myHash {
			<-blockchain
			log.Printf("%s: Consesus failed\nmax hash %s\nmy hash %s\n", host,
				consData.maxHash, consData.myHash)
			remotes := <-chRemotes
			chRemotes <- remotes
			requestFullBlockChain(remotes[0])
			if consData.onConsensusDone != nil {
				consData.onConsensusDone(false)
			}
		} else {
			log.Printf("%s: Consesus successful - Trans %s\n", host, vote)
			if consData.onConsensusDone != nil {
				consData.onConsensusDone(true)
			}

		}
	}

}
