package main

import (
	"bolt"
)

type BlockChainIterator struct {
	Db *bolt.DB
	//Hash指针
	CurrentHashPtr []byte
}

//func NewIterator(bc *BlockChain)  {
//
//}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		Db:bc.Db,
		//最初指向区块链的最后一个区块，随着Next的调用，不断变化
		CurrentHashPtr:bc.Tail,
	}
}


