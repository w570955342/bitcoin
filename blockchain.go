package main

import (
	"bolt"
	"log"
)

//4. 引入区块链
type BlockChain struct {
	//定义一个区块链数组
	//blocks []*Block
	Db   *bolt.DB //将区块链数据写入数据库
	Tail []byte  //存储最后一个区块的Hash
}

const BlockChainDb = "blockChain.db"
const BlockBucket = "blockBucket"

//5.创建区块链
func NewBlockChain(address string) *BlockChain {
	var lastBlockHash []byte
	//1. 打开数据库 每次测试程序需要删除 blockChain.db
	db, err := bolt.Open(BlockChainDb, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//defer db.Close()  //关闭后就不能持续添加区块了，除非添加的时候再次打开

	db.Update(func(tx *bolt.Tx) error {
		//2. 找到抽屉,没有就创建
		bucket := tx.Bucket([]byte(BlockBucket))
		if bucket == nil { //没有抽屉
			bucket, err = tx.CreateBucket([]byte(BlockBucket))
			if err != nil {
				log.Panic(err)
			}

			//创建一个创世块，并作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock(address)

			//3. 往bucket写数据
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("lastBlockHash"), genesisBlock.Hash)
			lastBlockHash=genesisBlock.Hash

			//测试数据，测试结束删除
			//blockBytes:=bucket.Get(genesisBlock.Hash)
			//block:=Deserialize(blockBytes)
			//fmt.Printf("刚写入的区块信息%x\n",block.Hash)

		}else {
			lastBlockHash=bucket.Get([]byte("lastBlockHash"))
		}
		return nil

	})

	return &BlockChain{
		Db:db,
		Tail:lastBlockHash,
	}
}

//创建一个创世区块
func GenesisBlock(address string) *Block {
	GenesisBlockCoinBaseTX:=NewCoinbaseTX(address,"这是一个关于创世区块的故事")
	return NewBlock([]*Transaction{GenesisBlockCoinBaseTX}, []byte{})
}

//6. 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	db:=bc.Db
	lastBlockHash:=bc.Tail//最后一个区块的Hash

	db.Update(func(tx *bolt.Tx) error {

		//添加数据
		bucket:=tx.Bucket([]byte(BlockBucket))
		if bucket == nil {
			log.Panic("bucket 不应该为空！")
		}

		//把新的区块写到数据库中 blockChain.db
		block := NewBlock(txs, lastBlockHash)
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("lastBlockHash"), block.Hash)
		//更新内存中数据
		bc.Tail=block.Hash
		return nil
	})
}