package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 10.0

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

//2. 生成一笔挖矿交易
func NewCoinbaseTX(address, data string) *Transaction {
	//挖矿交易的特点：
	//1. 只有一个input
	//2. 无需引用交易id 表名是挖矿交易
	//3. 无需引用index 设置为-1 任意设置的
	//矿工由于挖矿时无需指定签名，所以这个sig字段可以由矿工自由填写数据，一般是填写矿池的名字 例如	BTC.com
	input := TXInput{[]byte("挖矿交易Id为空"), -1, data}
	output := TXOutput{reward, address}

	//对于挖矿交易来说，只有一个input和一个output
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	//计算交易ID
	tx.SetHash()
	return &tx
}

//为交易绑定方法，判断该交易是否为挖矿交易
func (tx *Transaction) IsCoinbaseTX() bool {
	//1. 只有一个input
	if len(tx.TXInputs) == 1 {
		input := tx.TXInputs[0]
		//2. 无需引用交易id 表名是挖矿交易
		//3. 无需引用index 设置为-1 任意设置的
		if len(input.Txid) == 0 && input.Index == -1 {
			return true
		}
	}
	return false
}

//3. 生成一笔普通交易（找零）
func NewOrdinaryTX(from, to string, amount float64, bc *BlockChain) *Transaction {
	//1. 找到足够UTXO
	utxos, totalMoney := bc.FindEnoughUTXO(from, amount)
	if totalMoney < amount {
		fmt.Println("余额不足，交易失败！")
		return nil
	}

	//2. 创建交易输入，将UTXO逐一转成对应的input
	//map[string][]uint64
	var inputs []TXInput
	var outputs []TXOutput
	for TXid, indexSlice := range utxos {
		for _, i := range indexSlice {
			input := TXInput{[]byte(TXid), int64(i), from}
			inputs = append(inputs, input)
		}
	}

	//3. 创建交易输出
	output:=TXOutput{amount,to}
	outputs=append(outputs,output)

	//4. 找零
	if totalMoney>amount {
		outputs=append(outputs,TXOutput{totalMoney-amount,from})
	}

	//5. 生成交易
	tx:=Transaction{[]byte{},inputs,outputs}
	tx.SetHash()
	return &tx
}
