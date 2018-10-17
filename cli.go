package main

import (
	"os"
	"fmt"
)

//这是一个用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA     "add data to blockchain"
	printChain               "print all blockchain data" 
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
	case "addBlock":
		//3. 执行相应动作
		fmt.Printf("添加区块")
	case "printChain":
		fmt.Printf("打印区块")
	default:
		fmt.Printf("无效的命令，请检查!")
		fmt.Printf(Usage)
	}
}
