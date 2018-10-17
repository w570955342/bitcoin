package main

//1. 交易结构
type Transaction struct {
	TXId      []byte     //交易ID
	TXInputs  []TXInput  //交易输入切片
	TXOutputs []TXOutput //交易输出切片
}

//交易输入，钱的来源
type TXInput struct {
	Txid []byte//预消费的UTXO所在交易的交易ID，肯定不是当前交易的交易ID，至少是上一笔或者更早交易
	Index int64//预消费的UTXO在自己的交易中的索引
	Sig string//解锁脚本，用地址模拟
}

//交易输出，钱的去向
type TXOutput struct {
	Value float64//转账金额
	PubKeyHash string//所定脚本，用地址模拟
}

//2. 交易方法
//3. 创建挖矿交易
//4. 根据交易调整数据
