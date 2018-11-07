package main

import (
	"github.com/boltdb/bolt"
	"log"
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

//迭代器是属于区块链的
//Next方式是属于迭代器的
//1. 返回当前的区块
//2. 指针前移
func (it *BlockChainIterator) Next() *Block {
	var block Block
	it.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockBucket))
		if bucket == nil {
			log.Panic("迭代器遍历时bucket不应该为空，请检查!")
		}

		blockTmp := bucket.Get(it.CurrentHashPtr)
		//解码动作
		block = Deserialize(blockTmp)
		//游标哈希左移
		it.CurrentHashPtr = block.PrevHash

		return nil
	})

	return &block
}


