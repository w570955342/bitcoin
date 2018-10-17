package main

//7. 重构代码
func main() {
	blockChain:=NewBlockChain()
	cli := CLI{blockChain}
	cli.Run()
	//blockChain.AddBlock("这位大哥用两万个比特币买了一张pizza！")
	//blockChain.AddBlock("比特币暴跌，挖矿市场一片鬼哭狼嚎！")
	//
	////创建迭代器
	//it:=blockChain.NewIterator()
	//
	////通过迭代器返回数据库中的区块
	//for i:=1;;i++{
	//	block:=it.Next()
	//	fmt.Printf("================================= 当前区块高度：%d =================================\n",i)
	//	fmt.Printf("前区块哈希值： 	%x\n", block.PrevHash)
	//	fmt.Printf("当前区块哈希值：	%x\n", block.Hash)
	//	fmt.Printf("区块数据 :		%s\n", block.Data)
	//
	//	if len(block.PrevHash) == 0 {
	//		fmt.Println("区块链遍历结束！")
	//		break
	//	}
	//}
}