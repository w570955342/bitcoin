package main

//7. 重构代码
func main() {
	blockChain:=NewBlockChain("这是一个关于创世区块的故事")
	cli := CLI{blockChain}
	cli.Run()
	//blockChain.AddBlock("这位大哥用两万个比特币买了一张pizza！")
	//blockChain.AddBlock("比特币暴跌，挖矿市场一片鬼哭狼嚎！")
	//

}