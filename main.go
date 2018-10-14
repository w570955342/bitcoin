package bitcoin

import (
	"crypto/sha256"
	"fmt"
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
type BlockChain struct {
	//定义一个区块链数组
	blocks []*Block
}

//5.创建区块链
func NewBlockChain() *BlockChain {
	//创建一个创世块，并作为第一个区块添加到区块链中
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}

//创建一个创世区块
func GenesisBlock() *Block{
	return NewBlock("这是一个关于创世区块的故事",[]byte{})
}

//6. 添加区块
func (bc *BlockChain)AddBlock(data string)  {
	lastBlock:=bc.blocks[len(bc.blocks)-1]
	prevHash:=lastBlock.Hash

	block:=NewBlock(data,prevHash)
	bc.blocks=append(bc.blocks,block)
}
//7. 重构代码

func main() {
	blockChain:=NewBlockChain()
	blockChain.AddBlock("这位大哥用两万个比特币买了一张pizza！")
	blockChain.AddBlock("比特币暴跌，挖矿市场一片鬼哭狼嚎！")

	for idx,block:=range blockChain.blocks{
		fmt.Printf("================================= 当前区块高度：%d =================================\n",idx)
		fmt.Printf("前区块哈希值： 	%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值：	%x\n", block.Hash)
		fmt.Printf("区块数据 :	%s\n", block.Data)
		//fmt.Println("成功生成新的区块！")
	}

}