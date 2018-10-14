package main

import (
	"fmt"
)
//7. 重构代码
func main() {
	blockChain:=NewBlockChain()
	blockChain.AddBlock("这位大哥用两万个比特币买了一张pizza！")
	blockChain.AddBlock("比特币暴跌，挖矿市场一片鬼哭狼嚎！")

	for idx,block:=range blockChain.blocks{
		fmt.Printf("================================= 当前区块高度：%d =================================\n",idx)
		fmt.Printf("前区块哈希值： 	%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值：	%x\n", block.Hash)
		fmt.Printf("区块数据 :		%s\n", block.Data)
		//fmt.Println("成功生成新的区块！")
	}

}