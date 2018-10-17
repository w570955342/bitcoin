package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const reward = 12.5

//1. 交易结构
type Transaction struct {
	TXId      []byte     //交易ID
	TXInputs  []TXInput  //交易输入切片
	TXOutputs []TXOutput //交易输出切片
}

//交易输入，钱的来源
type TXInput struct {
	Txid  []byte //预消费的UTXO所在交易的交易ID，肯定不是当前交易的交易ID，至少是上一笔或者更早交易
	Index int64  //预消费的UTXO在自己的交易中的索引
	Sig   string //解锁脚本，用地址模拟
}

//交易输出，钱的去向
type TXOutput struct {
	Value      float64 //转账金额
	PubKeyHash string  //所定脚本，用地址模拟
}

//设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXId = hash[:]
}

//2. 交易方法
func NewCoinbaseTX(address, data string) *Transaction {
	//挖矿交易的特点：
	//1. 只有一个input
	//2. 无需引用交易id 表名是挖矿交易
	//3. 无需引用index 设置为-1 任意设置的
	//矿工由于挖矿时无需指定签名，所以这个sig字段可以由矿工自由填写数据，一般是填写矿池的名字 例如	BTC.com
	input := TXInput{[]byte("挖矿交易Id为空"), -1, data}
	output := TXOutput{reward, address}

	//对于挖矿交易来说，只有一个input和一个output
	tx:=Transaction{[]byte{},[]TXInput{input},[]TXOutput{output}}
	//计算交易ID
	tx.SetHash()
	return &tx
}

//3. 创建挖矿交易
//4. 根据交易调整数据
