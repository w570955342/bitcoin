package main

import (
	"bolt"
	"log"
	"fmt"
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

//找到指定地址的所有UTXO
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXO []TXOutput
	//定义一个map来保存消费过的output，key是这个output的交易id，value是这个交易中索引的切片
	//map[交易id][]int64
	spentOutputs := make(map[string][]int64)

	//创建迭代器
	it := bc.NewIterator()

	for {
		//1.遍历区块
		block := it.Next()

		//2. 遍历交易
		for _, tx := range block.Transactions {
			fmt.Printf("current TXId : %x\n", tx.TXId)

			//3. 遍历output，找到和自己相关的UTXO(在添加output之前检查一下是否已经消耗过)
			for i, output := range tx.TXOutputs {
				fmt.Printf("current index : %d\n", i)

				//这个output和我们目标的地址相同，满足条件，加到返回UTXO数组中
				if output.PubKeyHash == address {
					UTXO = append(UTXO, output)
				}
			}

			//4. 遍历input，找到自己花费过的UTXO的集合(把自己消耗过的标示出来)
			for _, input := range tx.TXInputs {
				//判断一下当前这个input和目标（李四）是否一致，如果相同，说明这个是李四消耗过的output,就加进来
				if input.Sig == address {
					//spentOutputs := make(map[string][]int64)
					//indexSlice := spentOutputs[string(input.Txid)]//定义一个空切片
					//indexSlice = append(indexSlice, input.Index)
					spentOutputs[string(input.Txid)]=append(spentOutputs[string(input.Txid)], input.Index)
					//map[2222] = []int64{0}
					//map[3333] = []int64{0, 1}
				//indexSlice 中的index可能会重复，来自不同交易信息
				}
			}

		}

		if len(block.PrevHash) == 0 {
			break
			fmt.Printf("区块遍历完成退出!")
		}
	}


	return UTXO
}