package main

import "fmt"

func (cli *CLI) AddBlock(data string) {
	cli.bc.AddBlock(data)
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
		fmt.Printf("区块数据：      %s\n", block.Data)

		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历结束！")
			break
		}
	}
}