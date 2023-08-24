package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Blockchain struct {
	NumBlocks    int
	GenesisBlock *Block
	TailBlock    *Block

	// Base de datos simulada para CRUD:
	Database SimulatedDatabase
}

func (chain *Blockchain) AssertBlockchain() {
	cnt := 0
	for block := chain.GenesisBlock; block != nil; block = block.NextBlock {
		if block.PrevBlock != nil && block.PrevBlock.Hash != block.PrevHash {
			panic("Assert blockchain: block hashes dont match")
		}
		cnt++
	}
	if cnt != chain.NumBlocks {
		panic("Assert blockchain: block cnt vs chain.numblocks doesnt match")
	}
}

func MakeBlockchainFromJson(jsonData []byte) Blockchain {
	var blocks []Block
	err := json.Unmarshal(jsonData, &blocks)
	if err != nil {
		return Blockchain{}
	}
	chain := Blockchain{}
	chain.NumBlocks = len(blocks)
	for i := 0; i < int(chain.NumBlocks); i++ {
		if i == 0 {
			chain.GenesisBlock = &blocks[i]
		} else {
			blocks[i].PrevBlock = &blocks[i-1]
			blocks[i-1].NextBlock = &blocks[i]
			chain.TailBlock = &blocks[i]
		}
		chain.HandleTransaction(blocks[i].Data)
	}
	return chain
}

func (chain *Blockchain) HandleTransaction(transaction Transaction) bool {
	chain.Database.EnsureInitialized()
	switch transaction.Cmd {

	case TransCmdCreateCandidate:
		return chain.Database.CreateCandidate(transaction.CandidateData)

	case TransCmdCreateVotante:
		return chain.Database.CreateVotante(transaction.VotanteData)

	case TransCmdUpdateCandidate:
		return chain.Database.UpdateCandidate(transaction.CandidateData)

	case TransCmdUpdateVotante:
		return chain.Database.UpdateVotante(transaction.VotanteData)

	default:
		panic("no valid transaction cmd")
		return false
	}
}

func (chain *Blockchain) AddTransaction(transaction Transaction) bool {
	if chain.HandleTransaction(transaction) {
		block := NewBlock(chain.TailBlock, transaction)
		if chain.GenesisBlock == nil {
			chain.GenesisBlock = block
		}
		chain.TailBlock = block
		chain.NumBlocks++
		/// DEBUG
		jsonData, _ := json.MarshalIndent(block, "", "\t")
		fmt.Println("Blockchain::AddTransaction successful")
		fmt.Println(string(jsonData))
		return true
	}
	jsonData, _ := json.MarshalIndent(transaction, "", "\t")
	fmt.Println("Blockchain::AddTransaction failed")
	fmt.Println(string(jsonData))
	return false
}

func (chain *Blockchain) GetBlockchainArrayRaw() []Block {
	blockArray := make([]Block, 0, chain.NumBlocks)
	for block := chain.GenesisBlock; block != nil; block = block.NextBlock {
		blockArray = append(blockArray, *block)
	}
	return blockArray
}

func (chain *Blockchain) GetBlockchainJson() []byte {
	jsonData, err := json.MarshalIndent(chain.GetBlockchainArrayRaw(), "", "\t")
	if err != nil {
		panic("error getting blockchain json")
		return []byte{}
	}
	return jsonData
}

func (chain *Blockchain) PrintBlockchain() {
	fmt.Println(string(chain.GetBlockchainJson()))
}

func main_blockchain() {
	oldTest()
}

func CreateTestBlockchain() *Blockchain {
	b := new(Blockchain)

	for i := 0; i < 3; i++ {
		f := Transaction{}
		f.Sender = fmt.Sprint(i + 1)
		if i%2 > 0 {
			f.Cmd = TransCmdCreateVotante
			vote := VOTANTE{}
			vote.ID = (i / 2) + 1
			vote.Nombre = f.Sender
			vote.Apellido = f.Sender
			vote.DNI = f.Sender
			vote.Candidato_Voto = f.Sender
			vote.Lugar_Votacion = f.Sender
			f.VotanteData = vote
		} else {
			f.Cmd = TransCmdCreateCandidate
			cand := CANDIDATE{}
			cand.ID = (i / 2) + 1
			cand.Nombre = f.Sender
			cand.Apellido = f.Sender
			cand.DNI = f.Sender
			cand.Numero_votacion = f.Sender
			f.CandidateData = cand
		}
		b.AddTransaction(f)
	}
	return b

}

func oldTest() {
	// TestCode
	fmt.Println("hello")
	b := CreateTestBlockchain()
	b.Database.PrintDatabase()
	// Testing copying blockchain from json
	// b.PrintBlockchain()
	b.AssertBlockchain()
	copied := MakeBlockchainFromJson(b.GetBlockchainJson())
	// copied.PrintBlockchain()
	copied.AssertBlockchain()
	if !bytes.Equal(copied.GetBlockchainJson(), b.GetBlockchainJson()) {
		panic("copy doesnt match original")
	}

	// Testing copying block from json

	tailData, _ := json.MarshalIndent(*(b.TailBlock), "", "\t")
	copiedBlock := new(Block)
	// fmt.Println(string(tailData))
	_ = json.Unmarshal(tailData, copiedBlock)
	newData, _ := json.MarshalIndent(*copiedBlock, "", "\t")
	// fmt.Println(string(newData))
	if !bytes.Equal(newData, tailData) {
		panic("copied block doesnt match")
	}

}
