package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"crypto/ecdsa"
	"crypto/rand"
	"strings"
	"math/big"
	"crypto/elliptic"
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
	bc.SignTransaction(&tx,privateKey)
	return &tx
}

//为普通交易绑定Sign方法
//参数为：私钥，inputs里面引用的所有的交易实体 map[string]Transaction
//map[交易ID]Transaction
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//1. 创建一个当前交易的副本：txCopy，使用方法： TrimmedCopy：要把Signature和PubKey字段设置为nil
	txCopy := tx.TrimmedCopy()
	//2. 循环遍历txCopy的inputs，得到这个input索引的output的公钥哈希
	for i, input := range txCopy.TXInputs {
		prevTX := prevTXs[string(input.Txid)]
		if len(prevTX.TXId) == 0 {
			log.Panic("引用的交易无效")
		}

		//不要对input进行赋值，这是一个副本，为了找到公钥哈希，要对txCopy.TXInputs[xx]进行操作，否则无法把pubKeyHash传进txCopy.TXInputs[xx]
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash

		//所需要的三个数据都具备了，开始做哈希处理
		//3. 生成要签名的数据。要签名的数据一定是哈希值
		//a. 我们对每一个input都要签名一次，签名的数据是由当前input引用的output的哈希+当前的outputs（都承载在当前这个txCopy里面）
		//b. 要对这个拼好的txCopy进行哈希处理，SetHash得到TXID，这个TXID就是我们要签名最终数据。
		txCopy.SetHash()

		//还原，以免影响后面input的签名
		txCopy.TXInputs[i].PubKey = nil
		signDataHash := txCopy.TXId
		//4. 执行签名动作得到r,s字节流
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		if err != nil {
			log.Panic(err)
		}

		//5. 放到我们所签名的input的Signature中
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TXInputs[i].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	for _, input := range tx.TXInputs {
		inputs = append(inputs, TXInput{input.Txid, input.Index, nil, nil})
	}

	return Transaction{tx.TXId, inputs, tx.TXOutputs}
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.TXId))

	for i, input := range tx.TXInputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.TXOutputs{
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}


//分析校验：
//所需要的数据：公钥，数据(txCopy，生成哈希), 签名
//我们要对每一个签名过得input进行校验

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbaseTX() {
		return true
	}

	//1. 得到签名的数据
	txCopy := tx.TrimmedCopy()

	for i, input := range tx.TXInputs {
		prevTX := prevTXs[string(input.Txid)]
		if len(prevTX.TXId) == 0 {
			log.Panic("引用的交易无效")
		}

		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		txCopy.SetHash()
		dataHash := txCopy.TXId
		//2. 得到Signature, 反推会r,s
		signature := input.Signature //拆，r,s

		//a. 定义两个辅助的big.int
		r := big.Int{}
		s := big.Int{}

		//b. 拆分我们signature，平均分，前半部分给r, 后半部分给s
		r.SetBytes(signature[:len(signature)/2 ])
		s.SetBytes(signature[len(signature)/2:])

		//3. 拆解PubKey, X, Y 得到原生公钥
		pubKey := input.PubKey //拆，X, Y

		//a. 定义两个辅助的big.int
		X := big.Int{}
		Y := big.Int{}

		//b. pubKey，平均分，前半部分给X, 后半部分给Y
		X.SetBytes(pubKey[:len(pubKey)/2 ])
		Y.SetBytes(pubKey[len(pubKey)/2:])

		pubKeyOrigin := ecdsa.PublicKey{elliptic.P256(), &X, &Y}

		//4. Verify
		if !ecdsa.Verify(&pubKeyOrigin, dataHash, &r, &s) {
			return false
		}
	}

	return true
}