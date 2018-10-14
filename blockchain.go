package main

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