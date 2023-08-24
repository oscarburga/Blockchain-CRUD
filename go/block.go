package main

import "encoding/json"

type Block struct {
	Data      Transaction `json:"Data"`
	Hash      string      `json:"Hash"`
	PrevHash  string      `json:"PrevHash"`
	PrevBlock *Block      `json:"-"`
	NextBlock *Block      `json:"-"`
}

func NewBlock(prevBlock *Block, trans Transaction) *Block {
	block := new(Block)
	block.Data = trans
	transBytes, err := json.Marshal(trans)
	if err != nil {
		panic("error")
		return nil
	}
	block.Hash = GetHashOfBytes(transBytes)

	if prevBlock != nil {
		block.PrevHash = prevBlock.Hash
		block.PrevBlock = prevBlock
		if prevBlock.NextBlock != nil {
			panic("prevBlock.NextBlock was not nil")
		}
		prevBlock.NextBlock = block
	} else {
		block.PrevHash = ""
		block.PrevBlock = nil
	}
	return block
}

func NewBlockFromJson(jsonData []byte) *Block {
	block := new(Block)
	err := json.Unmarshal(jsonData, block)
	if err != nil {
		return nil
	}
	return block
}
