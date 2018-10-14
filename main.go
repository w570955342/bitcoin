package main

import (
	"fmt"
	"crypto/sha256"
)

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

//4. 引入区块链
//5. 添加区块
//6. 重构代码

func main() {
	block := NewBlock("用两万个比特币买了一张pizza！", []byte{})

	fmt.Printf("前区块哈希值： %x\n", block.PrevHash)
	fmt.Printf("当前区块哈希值： %x\n", block.Hash)
	fmt.Printf("区块数据 :%s\n", block.Data)
	fmt.Println("成功生成新的区块！")
}
