package main

import "crypto/sha256"

//1. 定义区块结构
type Block struct {
	//1. 前区块哈希
	PrevHash []byte
	//2. 当前区块哈希
	Hash []byte
	//3. 交易数据
	Data []byte
}

//2. 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{}, //先填空，后面再计算 //TODO
		Data:     []byte(data),
	}

	block.SetHash()
	return &block
}

//3. 为Block绑定方法生成哈希
func (block *Block) SetHash() {
	//1. 拼接数据
	blockInfo := append(block.PrevHash, block.Data...)

	//2. sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}