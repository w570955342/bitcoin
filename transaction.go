package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"crypto/ecdsa"
)

const reward = 50.0

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
	//Sig   string //解锁脚本，用地址模拟

	Signature []byte //真正的数字签名，由r，s拼成的[]byte

	//约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）
	//注意，是公钥，不是哈希，也不是地址
	PubKey []byte
}


//交易输出，钱的去向
type TXOutput struct {
	Value float64 //转账金额
	//PubKeyHash string  //所定脚本，用地址模拟

	PubKeyHash []byte //收款方的公钥的哈希，注意，是哈希而不是公钥，也不是地址
}

//TXOutput存储的字段是地址对应的的公钥哈希，需要由地址反推出公钥哈希，然后创建TXOutput
//为了能够得到公钥哈希，为TXOutput绑定一个方法
func (output *TXOutput) SetPubKeyHash(address string) {
	////1. 解码
	////2. 截取出公钥哈希：去除version（1字节），去除校验码（4字节）
	//addressByte := base58.Decode(address) //25字节
	//len := len(addressByte)
	//
	//pubKeyHash := addressByte[1:len-4]

	//真正的锁定动作！！！！！
	output.PubKeyHash = GetPubKeyHashFromAddress(address)
}

//给TXOutput提供一个创建的方法，否则无法调用SetPubKeyHash方法
func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}

	output.SetPubKeyHash(address)
	return &output
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
	//签名信息先填为空，创建完整交易后，最后做一次签名
	input := TXInput{[]byte{}, -1, nil,[]byte(data)}
	//output := TXOutput{reward, address}

	output:=NewTXOutput(reward,address)

	//对于挖矿交易来说，只有一个input和一个output
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
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

	//a. 创建交易之后要进行数字签名->所以需要私钥->把 wallet.dat 中的密钥对加载到内存中
	wallet := NewWallet()

	//b. 找到自己的密钥对，把地址作为key找到value（即密钥对）
	key := wallet.WalletMap[from]
	if key == nil {
		fmt.Printf("没有找到地址\"%s\"对应的密钥对，交易创建失败!\n",from)
		return nil
	}

	//c. 得到对应的公钥，私钥
	pubKey := key.PubKey
	privateKey := key.Private  //稍后再用

	//传递公钥的哈希，而不是传递地址
	pubKeyHash := HashPubKey(pubKey)

	//1. 找到足够UTXO
	utxos, totalMoney := bc.FindEnoughUTXO(pubKeyHash, amount)
	if totalMoney < amount {
		fmt.Printf("\"%s\"只有 %f 比特币，余额不足，交易失败！", from, totalMoney)
		return nil
	}

	//2. 创建交易输入，将UTXO逐一转成对应的input
	//map[string][]uint64
	var inputs []TXInput
	var outputs []TXOutput
	for TXid, indexSlice := range utxos {
		for _, i := range indexSlice {
			input := TXInput{[]byte(TXid), int64(i),nil,pubKey}
			inputs = append(inputs, input)
		}
	}

	//3. 创建交易输出
	//output := TXOutput{amount, to}
	output :=NewTXOutput(amount,to)
	outputs = append(outputs, *output)

	//4. 找零
	if totalMoney > amount {
		output :=NewTXOutput(totalMoney - amount,from)
		outputs = append(outputs, *output)
	}

	//5. 生成交易
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()

	//6. 签名
	prevTXs:=make(map[string]Transaction)

	tx.Sign(*privateKey,prevTXs)
	return &tx
}

//为普通交易绑定Sign方法
//参数为：私钥，inputs里面引用的所有的交易实体 map[string]Transaction
//map[交易ID]Transaction
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//具体签名功能
	//TODO
}
