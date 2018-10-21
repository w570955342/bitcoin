package main

import (
	"os"
	"fmt"
	"strconv"
)

//这是一个用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	bc *BlockChain
}

const Usage = `
	printChain			"print all blockchain data" 
	getBalance --address ADDRESS	"获取指定地址ADDRESS的余额"
	send FROM TO AMOUNT MINER DATA	"由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
	newKey				"ecdsa P256单独创建一个密钥对，不保存在wallet.dat中"
	newWallet			"创建一个新的钱包(私钥公钥对)"
	listAddresses			"列举所有的密钥对地址"
`

//接受参数的动作，我们放到一个函数中

func (cli *CLI) Run() {

	//./block printChain
	//./block addBlock
	//1. 得到所有的命令
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}

	//2. 分析命令
	cmd := args[1]
	switch cmd {
	case "printChain":
		fmt.Printf("打印区块...\n")
		cli.PrinBlockChain()
	case "getBalance":
		fmt.Printf("获取余额...\n")
		if len(args) == 4 && args[2] == "--address" {
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Printf("转账开始...\n")
		if len(args) != 7 {
			fmt.Printf("参数个数错误，请检查！\n")
			fmt.Printf(Usage)
			return
		}
		//./block send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64) //知识点，请注意
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	case "newKey":
		fmt.Printf("ecdsa P256创建新的秘钥对...\n")
		cli.NewKey()
	case "newWallet":
		fmt.Printf("创建新的钱包...\n")
		cli.NewWallet()
	case "listAddresses":
		fmt.Printf("列举所有地址...\n")
		cli.ListAddresses()
	default:
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Usage)
	}
}
