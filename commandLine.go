package main

import "fmt"

func (cli *CLI) AddBlock(data string) {
	//cli.bc.AddBlock(txs) //todo
	fmt.Printf("添加区块成功！\n")
}

func (cli *CLI) PrinBlockChain() {
	bc := cli.bc
	//创建迭代器
	it:=bc.NewIterator()

	//通过迭代器返回数据库中的区块
	for i:=1;;i++{
		block:=it.Next()
		fmt.Printf("================================= 当前区块高度：%d =================================\n",i)
		fmt.Printf("前区块哈希值：  %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值：%x\n", block.Hash)
		fmt.Printf("区块数据：      %s\n", block.Transactions[0].TXInputs[0].Sig)

		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历结束！")
			break
		}
	}
}

func (cli *CLI) GetBalance(address string) {

	utxos := cli.bc.FindUTXOs(address)

	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("\"%s\"的余额为：%f\n", address, total)
}

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	fmt.Printf("from : %s\n", from)
	fmt.Printf("to : %s\n", to)
	fmt.Printf("amount : %f\n", amount)
	fmt.Printf("miner : %s\n", miner)
	fmt.Printf("data : %s\n", data)

	//具体的逻辑，TODO
}
